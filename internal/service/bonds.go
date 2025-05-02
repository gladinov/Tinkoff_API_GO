package service

type Bond struct {
	Identifiers              Identifiers
	Name                     string  // GetBondsActionsFromPortfolio
	InstrumentType           string  // T_Api_Getportfolio
	Currency                 string  // T_Api_Getportfolio
	Quantity                 float64 // T_Api_Getportfolio
	AveragePositionPrice     float64 // T_Api_Getportfolio
	ExpectedYield            float64 // T_Api_Getportfolio
	CurrentNkd               float64 // T_Api_Getportfolio
	CurrentPrice             float64 // T_Api_Getportfolio
	AveragePositionPriceFifo float64 // T_Api_Getportfolio
	Blocked                  bool    // T_Api_Getportfolio
	ExpectedYieldFifo        float64 // T_Api_Getportfolio
	DailyYield               float64 // T_Api_Getportfolio
	// Amortizations            *moex.Amortizations // GetBondsActionsFromPortfolio
	// Coupons                  *moex.Coupons       // GetBondsActionsFromPortfolio
	// Offers                   *moex.Offers        // GetBondsActionsFromPortfolio
	// Duration                 moex.Duration       // GetBondsActionsFromPortfolio
}

type Identifiers struct {
	Ticker        string // GetBondsActionsFromPortfolio
	ClassCode     string // GetBondsActionsFromPortfolio
	Figi          string // T_Api_Getportfolio
	InstrumentUid string // T_Api_Getportfolio
	PositionUid   string // T_Api_Getportfolio
	AssetUid      string // GetBondsActionsFromPortfolio
}

// Получение данных с московской биржи
// func (b *Bond) GetActionFromMoex() error {
// 	MoexUnmarshallData := moex.MoexUnmarshallStruct{}
// 	err := MoexUnmarshallData.GetBondsFromMOEX(b.Identifiers.Ticker, 0, 20)
// 	if err != nil {
// 		return errors.New("GetBondsActionsFromPortfolio: GetBondsFromMOEX" + err.Error())
// 	}
// 	err = MoexUnmarshallData.GetDurationFromMoex(b.Identifiers.Ticker, b.Identifiers.ClassCode)
// 	if err != nil {
// 		return errors.New("GetBondsActionsFromPortfolio: GetBondsFromMOEX" + err.Error())
// 	}
// 	b.Amortizations = MoexUnmarshallData.Amortizations
// 	b.Offers = MoexUnmarshallData.Offers
// 	b.Coupons = MoexUnmarshallData.Coupons
// 	b.Duration = MoexUnmarshallData.Duration
// 	return nil
// }
