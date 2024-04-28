package entity

import "errors"

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, price, tax, finalPrice float64) (*Order, error) {
	order := &Order{ID: id, Price: price, Tax: tax, FinalPrice: finalPrice}
	err := order.IsValid()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) IsValid() error {
	if o.ID == "" {
		return errors.New("invalid ID")
	}
	if o.Price <= 0 {
		return errors.New("invalid Price")
	}
	if o.Tax <= 0 {
		return errors.New("invalid Tax")
	}
	return nil
}

func (o *Order) CalculateFinalPrice() error {
	err := o.IsValid()
	if err != nil {
		return err
	}
	o.FinalPrice = o.Price + o.Tax
	return nil
}
