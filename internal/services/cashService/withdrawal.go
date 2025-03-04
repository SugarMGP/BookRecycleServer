package cashService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/billService"
	"bookrecycle-server/pkg/database"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Withdrawal 提现
func Withdrawal(user *models.User, newBalance decimal.Decimal, cash decimal.Decimal) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		user.Balance = newBalance.StringFixedBank(2)
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		err := billService.SaveBill(&models.Bill{
			User:   user.ID,
			Amount: cash.Neg().StringFixedBank(2),
		})
		if err != nil {
			return err
		}

		return tx.Create(&models.Withdrawal{
			Money:   cash.StringFixedBank(2),
			Account: user.Phone,
			Name:    user.Name,
		}).Error
	})
}

// GetWithdrawalList 获取提现记录
func GetWithdrawalList() ([]*models.Withdrawal, error) {
	var list []*models.Withdrawal
	result := database.DB.Find(&list)
	return list, result.Error
}
