package cashService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/billService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/pkg/database"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Withdrawal 提现
func Withdrawal(userID uint, amount float64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		user, err := userService.GetUserByID(userID)
		if err != nil {
			return err
		}

		balance, _ := decimal.NewFromString(user.Balance)
		cash := decimal.NewFromFloat(amount)
		user.Balance = balance.Sub(cash).StringFixedBank(2)
		if err = tx.Save(user).Error; err != nil {
			return err
		}

		err = billService.SaveBill(&models.Bill{
			User:   userID,
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
