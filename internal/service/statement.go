package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type statementService struct {
	account     repository.AccountRepository
	transaction repository.TransactionRepository
	journal     repository.JournalEntryRepository
}

func NewStatementService(
	account repository.AccountRepository,
	transaction repository.TransactionRepository,
	journal repository.JournalEntryRepository,
) *statementService {
	return &statementService{
		account:     account,
		transaction: transaction,
		journal:     journal,
	}
}

func (s *statementService) GenerateStatement(ctx context.Context, accountID uuid.UUID, w io.Writer) error {

	// Step 1 — fetch account
	account, err := s.account.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("statementService: fetching account: %w", err)
	}

	// Step 2 — fetch transactions (last 100)
	transactions, err := s.transaction.GetByAccountID(ctx, accountID, 100, nil)
	if err != nil {
		return fmt.Errorf("statementService: fetching transactions: %w", err)
	}

	// Step 3 — fetch current balance
	balance, err := s.journal.GetBalance(ctx, accountID)
	if err != nil {
		return fmt.Errorf("statementService: fetching balance: %w", err)
	}

	// Step 4 — generate PDF
	// --- Design Tokens (Modern Fintech Palette) ---
	navy := []int{15, 23, 42}       // Deep Slate Navy
	slate := []int{100, 116, 139}   // Cool Grey
	emerald := []int{16, 185, 129}  // Success Green
	lightBg := []int{248, 250, 252} // Subtle background

	// --- PDF Initialization ---
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// 1. Header: Brand Bar
	pdf.SetFillColor(navy[0], navy[1], navy[2])
	pdf.Rect(0, 0, 210, 40, "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 22)
	pdf.Text(15, 22, "LEDGER-CORE")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(148, 163, 184)
	pdf.Text(15, 30, "Automated Financial Report")

	// 2. Account Summary Card
	pdf.SetY(45)
	pdf.SetFillColor(lightBg[0], lightBg[1], lightBg[2])
	pdf.Rect(15, 45, 180, 35, "F") // Summary Box

	// Left: Account Details
	pdf.SetTextColor(slate[0], slate[1], slate[2])
	pdf.SetFont("Arial", "B", 8)
	pdf.SetXY(20, 52)
	pdf.Cell(0, 0, "ACCOUNT NUMBER")

	pdf.SetTextColor(navy[0], navy[1], navy[2])
	pdf.SetFont("Arial", "B", 11)
	pdf.SetXY(20, 58)
	pdf.Cell(0, 0, account.AccountNumber)

	pdf.SetTextColor(slate[0], slate[1], slate[2])
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(20, 65)
	pdf.Cell(0, 0, fmt.Sprintf("ID: %s", account.ID.String()))

	// Right: Balance Highlight
	pdf.SetTextColor(slate[0], slate[1], slate[2])
	pdf.SetFont("Arial", "B", 8)
	pdf.SetXY(130, 52)
	pdf.CellFormat(60, 0, "AVAILABLE BALANCE", "", 0, "R", false, 0, "")

	pdf.SetTextColor(emerald[0], emerald[1], emerald[2])
	pdf.SetFont("Arial", "B", 18)
	pdf.SetXY(130, 62)
	pdf.CellFormat(60, 0, fmt.Sprintf("Rs. %.2f", float64(balance)/100), "", 0, "R", false, 0, "")

	// 3. Transaction Table
	pdf.SetY(90)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(navy[0], navy[1], navy[2])
	pdf.SetTextColor(255, 255, 255)

	// Custom Table Header (No borders, just a solid bar)
	pdf.CellFormat(35, 10, "  TX ID", "", 0, "L", true, 0, "")
	pdf.CellFormat(25, 10, "TYPE", "", 0, "C", true, 0, "")
	pdf.CellFormat(25, 10, "STATUS", "", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "DATE", "", 0, "C", true, 0, "")
	pdf.CellFormat(50, 10, "AMOUNT (INR)  ", "", 1, "R", true, 0, "")

	// 4. Transaction Rows
	pdf.SetFont("Arial", "", 9)
	for i, tx := range transactions {
		// Clean Zebra Striping
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255)
		} else {
			pdf.SetFillColor(252, 253, 255)
		}

		pdf.SetTextColor(navy[0], navy[1], navy[2])

		// Row content with "B" (Bottom border) only for a clean line look
		pdf.SetDrawColor(241, 245, 249)

		txID := "#" + tx.ID.String()[:7]
		date := tx.CreatedAt.Format("02 Jan 2006")

		pdf.CellFormat(35, 9, "  "+txID, "B", 0, "L", true, 0, "")
		pdf.CellFormat(25, 9, strings.ToUpper(tx.Type), "B", 0, "C", true, 0, "")

		// Semantic Status Coloring
		if tx.Status == "COMPLETED" {
			pdf.SetTextColor(emerald[0], emerald[1], emerald[2])
		} else {
			pdf.SetTextColor(239, 68, 68)
		}
		pdf.CellFormat(25, 9, tx.Status, "B", 0, "C", true, 0, "")

		pdf.SetTextColor(slate[0], slate[1], slate[2])
		pdf.CellFormat(45, 9, date, "B", 0, "C", true, 0, "")

		pdf.SetTextColor(navy[0], navy[1], navy[2])
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(50, 9, fmt.Sprintf("%.2f  ", float64(tx.Amount)/100), "B", 1, "R", true, 0, "")
		pdf.SetFont("Arial", "", 9)
	}

	// 5. Footer
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(148, 163, 184)
	pdf.CellFormat(0, 10, "ledger-core | System Generated Statement | "+time.Now().Format("Jan 02, 2006"), "T", 0, "C", false, 0, "")

	return pdf.Output(w)
}
