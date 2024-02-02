package bill

import "marketplace_server/internal/bill/model"

type BillAppInterface interface {
	CreateBill(bill *model.Transaction) error
}

var _ BillAppInterface = &BillApp{}

type BillApp struct {
	BillRepo BillRepo
}

func NewBillApp(billRepo BillRepo) *BillApp {
	return &BillApp{
		BillRepo: billRepo,
	}
}

func (a *BillApp) CreateBill(bill *model.Transaction) error {
	return a.BillRepo.Save(bill)
}