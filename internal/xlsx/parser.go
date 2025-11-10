package xlsx

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joelovien/go-xlsx-api/internal/models"
	"github.com/xuri/excelize/v2"
)

type Parser struct {
	workerPoolSize int
}

func NewParser(workerPoolSize int) *Parser {
	return &Parser{
		workerPoolSize: workerPoolSize,
	}
}

type ParseResult struct {
	UploadID     string
	Records      []models.Record
	RowsAccepted int
	RowsRejected int
	Errors       []string
}

func (p *Parser) Parse(ctx context.Context, reader io.Reader, uploadID string) (*ParseResult, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("xlsx file has no sheets")
	}

	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName, excelize.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet %s: %w", sheetName, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("xlsx file has no data")
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("xlsx file must have at least header row and one data row")
	}

	
	headerRowIndex := 0
	dataStartIndex := 1

	if len(rows[0]) == 1 && len(rows) > 8 {
		headerRowIndex = 7
		dataStartIndex = 8
	}

	if headerRowIndex >= len(rows) {
		return nil, fmt.Errorf("xlsx file does not have enough rows")
	}

	if len(rows[headerRowIndex]) == 0 {
		return nil, fmt.Errorf("xlsx file has no headers")
	}

	headers := make([]string, len(rows[headerRowIndex]))
	for i, cell := range rows[headerRowIndex] {
		headers[i] = strings.TrimSpace(cell)
	}

	hasHeader := false
	for _, header := range headers {
		if header != "" {
			hasHeader = true
			break
		}
	}

	if !hasHeader {
		return nil, fmt.Errorf("xlsx file has no valid headers")
	}

	dataRows := rows[dataStartIndex:]
	result := &ParseResult{
		UploadID: uploadID,
		Records:  make([]models.Record, 0, len(dataRows)),
		Errors:   make([]string, 0),
	}

	type rowJob struct {
		index int
		row   []string
	}

	jobs := make(chan rowJob, len(dataRows))
	results := make(chan models.ParsedRow, len(dataRows))

	// Start worker pool
	for w := 0; w < p.workerPoolSize; w++ {
		go func() {
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					parsed := p.parseRow(headers, job.row, job.index)
					results <- parsed
				}
			}
		}()
	}

	for i, row := range dataRows {
		normalizedRow := make([]string, len(headers))
		copy(normalizedRow, row)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case jobs <- rowJob{index: i, row: normalizedRow}:
		}
	}
	close(jobs)

	for i := 0; i < len(dataRows); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case parsed := <-results:
			if parsed.Valid {
				record := models.Record{
					ID:        uuid.New().String(),
					UploadID:  uploadID,
					Data:      parsed.Data,
					CreatedAt: result.CreatedAt(),
				}
				result.Records = append(result.Records, record)
				result.RowsAccepted++
			} else {
				result.RowsRejected++
				if parsed.Error != "" {
					result.Errors = append(result.Errors, parsed.Error)
				}
			}
		}
	}

	return result, nil
}

func (p *Parser) parseRow(headers []string, row []string, transactionIndex int) models.ParsedRow {
	// Skip completely empty rows
	if p.isEmptyRow(row) {
		return models.ParsedRow{
			Valid: false,
			Error: "empty row",
		}
	}

	data := make(map[string]interface{})

	data["Transaction Index"] = transactionIndex + 1

	for i, header := range headers {
		var value interface{}
		if i < len(row) {
			cellValue := strings.TrimSpace(row[i])
			if cellValue != "" {
				value = cellValue
			}
		}
		// Skip empty header names
		if strings.TrimSpace(header) != "" {
			data[header] = value
		}
	}

	return models.ParsedRow{
		Data:  data,
		Valid: true,
	}
}

func (p *Parser) isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func (pr *ParseResult) CreatedAt() time.Time {
	return time.Now()
}
