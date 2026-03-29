package service

import (
	"math"
	"sort"

	repository "github.com/harishash/dotshop-be/internal/repositories"

	"time"

	"github.com/harishash/dotshop-be/internal/dto"
)

const (
	CURATOR_COMISSION_RATE = 0.15
)

type CuratorDashboardService interface {
	GraphDataForSales(curatorID uint, startDate, endDate *time.Time) (*dto.SalesGraphResponse, error)
	GraphDataForRevenue(curatorID uint, startDate, endDate *time.Time) (*dto.RevenueGraphResponse, error)
	GraphDataForOrder(curatorID uint, startDate, endDate *time.Time) (*dto.OrderGraphResponse, error)
	GraphDataForReturns(curatorID uint, startDate, endDate *time.Time) (*dto.ReturnsGraphResponse, error)
	GraphDataForAvgOrderValue(curatorID uint, startDate, endDate *time.Time) (*dto.AvgOrderValueGraphResponse, error)
	GraphDataForAUPOrder(curatorID uint, startDate, endDate *time.Time) (*dto.AvgUnitsPerOrderGraphResponse, error)
	GraphDataForUnitsSold(curatorID uint, startDate, endDate *time.Time) (*dto.UnitsSoldGraphResponse, error)
	GetCuratorTopWishlist(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error)
	GetCuratorTopSellingProducts(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error)
	GetCuratorTopSellingBrands(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error)
	GetCuratorTopPurchasers(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error)
	GetCuratorSaleByCategory(curatorId uint, query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error)
}

type curatorDashboardService struct {
	repo repository.CuratorDashboardRepository
}

func NewCuratorDashboardService(repo repository.CuratorDashboardRepository) CuratorDashboardService {
	return &curatorDashboardService{repo}
}

func (s *curatorDashboardService) GraphDataForOrder(curatorID uint, startDate, endDate *time.Time) (*dto.OrderGraphResponse, error) {
	startDateStr, endDateStr, interval := intervalSelection(startDate, endDate)
	orders, err := s.repo.GetOrders(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	getIntervalStart := func(order dto.OrderCount) time.Time {
		return order.IntervalStart
	}

	getCount := func(order dto.OrderCount) uint {
		return order.OrderCount
	}

	createOutput := func(intervalStart time.Time, count uint) dto.OrderCount {
		return dto.OrderCount{
			IntervalStart: intervalStart,
			OrderCount:    count,
		}
	}

	getOutputIntervalStart := func(order dto.OrderCount) time.Time {
		return order.IntervalStart
	}

	data := aggregateByInterval(orders, interval, getIntervalStart, getCount, createOutput, getOutputIntervalStart, func(a, b uint) uint { return a + b })

	var totalOrders uint
	for _, order := range data {
		totalOrders += uint(order.OrderCount)
	}

	orderForGraph := &dto.OrderGraphResponse{
		TotalOrderCount: totalOrders,
		Data:            data,
	}

	return orderForGraph, nil
}

func aggregateByInterval[T any, V any, U any](data []T, interval string, getIntervalStart func(T) time.Time, getCount func(T) V, createOutput func(time.Time, V) U, getOutputIntervalStart func(U) time.Time, add func(V, V) V) []U {
	var aggregatedOrders []U
	orderMap := make(map[string]V)

	layout := "2006-01-02"
	switch interval {
	case "month":
		layout = "2006-01"
	case "year":
		layout = "2006"
	}

	for _, item := range data {
		intervalStart := getIntervalStart(item).Format(layout)
		orderMap[intervalStart] = add(orderMap[intervalStart], getCount(item))
	}

	for intervalStart, count := range orderMap {
		parsedTime, _ := time.Parse(layout, intervalStart)
		aggregatedOrders = append(aggregatedOrders, createOutput(parsedTime, count))
	}

	sort.Slice(aggregatedOrders, func(i, j int) bool {
		return getOutputIntervalStart(aggregatedOrders[i]).Before(getOutputIntervalStart(aggregatedOrders[j]))
	})

	return aggregatedOrders
}

func (s *curatorDashboardService) GraphDataForReturns(curatorID uint, startDate, endDate *time.Time) (*dto.ReturnsGraphResponse, error) {
	// Determine date range and interval
	startDateStr, endDateStr, interval := intervalSelection(startDate, endDate)

	// Fetch return data from repository
	returns, err := s.repo.GetOrderReturns(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	// Fetch units sold data
	unitsSold, err := s.GraphDataForUnitsSold(curatorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	getIntervalStart := func(returns dto.OrderReturns) time.Time {
		return returns.IntervalStart
	}

	getCount := func(returns dto.OrderReturns) uint {
		return returns.TotalReturnedQuantity
	}

	createOutput := func(intervalStart time.Time, count uint) dto.OrderReturns {
		return dto.OrderReturns{
			IntervalStart:         intervalStart,
			TotalReturnedQuantity: count,
		}
	}

	getOutputIntervalStart := func(returns dto.OrderReturns) time.Time {
		return returns.IntervalStart
	}

	aggregatedReturns := aggregateByInterval(returns, interval, getIntervalStart, getCount, createOutput, getOutputIntervalStart, func(a, b uint) uint { return a + b })

	// Calculate total returns
	var totalReturns uint
	for _, returnOrder := range aggregatedReturns {
		totalReturns += uint(returnOrder.TotalReturnedQuantity)
	}

	// Calculate return rate with check for division by zero
	returnRate := 0.0
	if unitsSold.TotalUnitsSold > 0 {
		returnRate = (float64(totalReturns) / float64(unitsSold.TotalUnitsSold)) * 100
	}

	// Prepare response DTO
	returnsForGraph := &dto.ReturnsGraphResponse{
		TotalReturns: totalReturns,
		ReturnRate:   returnRate,
		Data:         aggregatedReturns,
	}

	return returnsForGraph, nil
}

func (s *curatorDashboardService) GraphDataForSales(
	curatorID uint,
	startDate, endDate *time.Time) (
	*dto.SalesGraphResponse, error) {

	startDateStr, endDateStr, interval := intervalSelection(startDate, endDate)

	sales, err := s.repo.GetSales(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	getIntervalStart := func(sales dto.SalesIntervalResult) time.Time {
		return sales.IntervalStart
	}

	getCount := func(sales dto.SalesIntervalResult) float64 {
		return sales.TotalAmountSum
	}

	createOutput := func(intervalStart time.Time, count float64) dto.SalesIntervalResult {
		return dto.SalesIntervalResult{
			IntervalStart:  intervalStart,
			TotalAmountSum: count,
		}
	}

	getOutputIntervalStart := func(sales dto.SalesIntervalResult) time.Time {
		return sales.IntervalStart
	}

	data := aggregateByInterval(sales, interval, getIntervalStart, getCount, createOutput, getOutputIntervalStart, func(a, b float64) float64 { return a + b })

	var totalSales float64
	for _, order := range data {
		totalSales += order.TotalAmountSum
	}

	salesForGraph := &dto.SalesGraphResponse{
		TotalSales: totalSales,
		Data:       data,
	}
	return salesForGraph, nil
}

func (s *curatorDashboardService) GraphDataForRevenue(
	curatorID uint,
	startDate, endDate *time.Time) (
	*dto.RevenueGraphResponse,
	error) {

	startDateStr, endDateStr, interval := intervalSelection(startDate, endDate)

	revenues, err := s.repo.GetRevenue(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	getIntervalStart := func(revenu dto.RevenueIntervalResult) time.Time {
		return revenu.IntervalStart
	}

	getCount := func(revenu dto.RevenueIntervalResult) float64 {
		return revenu.TotalRevenue
	}

	createOutput := func(intervalStart time.Time, count float64) dto.RevenueIntervalResult {
		return dto.RevenueIntervalResult{
			IntervalStart: intervalStart,
			TotalRevenue:  count,
		}
	}

	getOutputIntervalStart := func(revenu dto.RevenueIntervalResult) time.Time {
		return revenu.IntervalStart
	}

	data := aggregateByInterval(revenues, interval, getIntervalStart, getCount, createOutput, getOutputIntervalStart, func(a, b float64) float64 { return a + b })

	var totalRevenue float64
	for _, revenue := range data {
		totalRevenue += revenue.TotalRevenue
	}

	for i := range data {
		data[i].TotalRevenue *= CURATOR_COMISSION_RATE
	}

	revenueForGraph := &dto.RevenueGraphResponse{
		TotalRevenue: totalRevenue * CURATOR_COMISSION_RATE,
		Data:         data,
	}
	return revenueForGraph, nil

}

func (s *curatorDashboardService) GraphDataForAvgOrderValue(
	curatorID uint,
	startDate, endDate *time.Time) (
	*dto.AvgOrderValueGraphResponse, error) {

	startDateStr, endDateStr, _ := intervalSelection(startDate, endDate)

	// Get average order value for a day
	avgOrderValue, err := s.repo.GetAverageOrderValue(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	// data := aggregateByInterval(avgOrderValue, interval, getIntervalStart, getCount, createOutput, getOutputIntervalStart, func(a, b float64) float64 { return a + b })
	data, _ := aggregateByMonth(*startDate, *endDate, avgOrderValue)

	responseData := []dto.AOVIntervalResult{}
	var totalOrderValue float64
	var numberOfOrders float64

	for _, value := range data {
		if value.TotalNumberOfOrders == 0 {
			responseData = append(responseData, dto.AOVIntervalResult{
				IntervalStart:     value.IntervalStart,
				AverageOrderValue: 0,
			})
			continue
		}
		responseData = append(responseData, dto.AOVIntervalResult{
			IntervalStart:     value.IntervalStart,
			AverageOrderValue: value.TotalOrderValue / value.TotalNumberOfOrders,
		})
		numberOfOrders += value.TotalNumberOfOrders
	}

	for _, value := range responseData {
		totalOrderValue += value.AverageOrderValue
	}

	avgOrderValueGraphResponse := &dto.AvgOrderValueGraphResponse{
		AvgOrderValue: totalOrderValue / numberOfOrders,
		Data:          responseData,
	}
	return avgOrderValueGraphResponse, nil
}

func (s *curatorDashboardService) GraphDataForAUPOrder(
	curatorID uint,
	startDate, endDate *time.Time) (
	*dto.AvgUnitsPerOrderGraphResponse, error) {

	startDateStr, endDateStr, _ := intervalSelection(startDate, endDate)

	avgUnitsPerOrder, err := s.repo.GetAverageUnitsPerOrder(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	data, _ := aggregateByMonthForUnitsPerOrder(*startDate, *endDate, avgUnitsPerOrder)

	responseData := []dto.AUPIntervalResult{}
	var totalUnitsSold float64
	var numberOfOrders float64

	for _, value := range data {
		if value.TotalOrders == 0 {
			responseData = append(responseData, dto.AUPIntervalResult{
				IntervalStart:        value.IntervalStart,
				AverageUnitsPerOrder: 0,
			})
			continue
		}
		responseData = append(responseData, dto.AUPIntervalResult{
			IntervalStart:        value.IntervalStart,
			AverageUnitsPerOrder: value.TotalUnitsSold / value.TotalOrders,
		})
		numberOfOrders++
	}
	for _, value := range responseData {
		totalUnitsSold += value.AverageUnitsPerOrder
	}
	avgOrderValueGraphResponse := &dto.AvgUnitsPerOrderGraphResponse{
		AvgUnitsPerOrder: totalUnitsSold / numberOfOrders,
		Data:             responseData,
	}
	return avgOrderValueGraphResponse, nil
}

func (s *curatorDashboardService) GraphDataForUnitsSold(
	curatorID uint,
	startDate, endDate *time.Time) (
	*dto.UnitsSoldGraphResponse,
	error) {

	startDateStr, endDateStr, interval := intervalSelection(startDate, endDate)

	units, err := s.repo.GetUnitsSold(startDateStr, endDateStr, "day", curatorID)
	if err != nil {
		return nil, err
	}

	getIntervalStart := func(order dto.UnitsSold) time.Time {
		return order.IntervalStart
	}

	getCount := func(order dto.UnitsSold) uint {
		return order.UnitsSold
	}

	createOutput := func(intervalStart time.Time, count uint) dto.UnitsSold {
		return dto.UnitsSold{
			IntervalStart: intervalStart,
			UnitsSold:     count,
		}
	}

	getOutputIntervalStart := func(order dto.UnitsSold) time.Time {
		return order.IntervalStart
	}

	data := aggregateByInterval(units, interval, getIntervalStart, getCount, createOutput, getOutputIntervalStart, func(a, b uint) uint { return a + b })

	var totalUnitsSold uint
	for _, unit := range data {
		totalUnitsSold += unit.UnitsSold
	}

	unitsSoldGraphResponse := &dto.UnitsSoldGraphResponse{
		TotalUnitsSold: totalUnitsSold,
		Data:           data,
	}
	return unitsSoldGraphResponse, nil
}

func (s *curatorDashboardService) GetCuratorTopWishlist(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.repo.GetCuratorTopWishlistProduct(curatorId, query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}
	return res, nil
}

func (s *curatorDashboardService) GetCuratorTopSellingProducts(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.repo.GetCuratorTopSellingProducts(curatorId, query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}
	return res, nil
}

func (s *curatorDashboardService) GetCuratorTopSellingBrands(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.repo.GetCuratorTopSellingBrands(curatorId, query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}
	return res, nil
}

func (s *curatorDashboardService) GetCuratorTopPurchasers(curatorId uint, query *dto.CommonProductRequest) (*dto.Response, error) {
	data, paging, err := s.repo.GetCuratorTopPurchasers(curatorId, query)
	if err != nil {
		return nil, err
	}

	res := &dto.Response{
		Data:   data,
		Paging: *paging,
	}
	return res, nil
}

func (s *curatorDashboardService) GetCuratorSaleByCategory(curatorId uint, query *dto.SaleRequest) ([]dto.SaleByCategoryResponse, error) {
	return s.repo.GetCuratorSalesByCategory(curatorId, query)
}

func intervalSelection(startDate, endDate *time.Time) (string, string, string) {

	var interval string

	startDateStr := startDate.Format("2006-01-02T15:04:05Z")
	endDateStr := endDate.Format("2006-01-02T15:04:05Z")

	duration := endDate.Sub(*startDate)

	switch {
	case duration.Hours() <= 24:
		interval = "hour"
	case duration.Hours() <= 32*24:
		interval = "day"
	case duration.Hours() > 32*24:
		interval = "month"
	default:
		interval = "day"
	}

	return startDateStr, endDateStr, interval
}

func aggregateByMonth(from, to time.Time, data []dto.AOVIntervalResultResponse) ([]dto.AOVIntervalResultResponse, error) {

	// Calculate the difference between the 'from' and 'to' dates
	daysDiff := to.Sub(from).Hours() / 24

	// If the difference is greater than 31 days, aggregate by month
	if daysDiff > 31 {
		monthlyData := make(map[string]dto.AOVIntervalResultResponse)

		for _, record := range data {
			// Parse IntervalStart to time.Time for monthly aggregation
			intervalTime := record.IntervalStart

			// Format the date to "YYYY-MM" for monthly aggregation
			monthKey := intervalTime.Format("2006-01")

			if existingRecord, exists := monthlyData[monthKey]; exists {
				// Aggregate the values
				existingRecord.TotalOrderValue += record.TotalOrderValue
				existingRecord.TotalNumberOfOrders += record.TotalNumberOfOrders
				monthlyData[monthKey] = existingRecord
			} else {
				// Create a new entry in the map
				monthlyData[monthKey] = dto.AOVIntervalResultResponse{
					IntervalStart:       intervalTime, // ISO 8601 format
					TotalOrderValue:     record.TotalOrderValue,
					TotalNumberOfOrders: record.TotalNumberOfOrders,
				}
			}
		}

		// Convert the map to a slice of dto.AOVIntervalResultResponse
		result := make([]dto.AOVIntervalResultResponse, 0, len(monthlyData))
		for _, record := range monthlyData {
			// Handle NaN values by replacing with zero
			if math.IsNaN(record.TotalOrderValue) {
				record.TotalOrderValue = 0
			}
			if math.IsNaN(record.TotalNumberOfOrders) {
				record.TotalNumberOfOrders = 0
			}
			result = append(result, record)
		}
		return result, nil
	}

	// If the difference is 31 days or less, return the original data
	return data, nil
}

func aggregateByMonthForUnitsPerOrder(from, to time.Time, data []dto.UnitsSoldPerOrder) ([]dto.UnitsSoldPerOrder, error) {

	// Calculate the difference between the 'from' and 'to' dates
	daysDiff := to.Sub(from).Hours() / 24

	// If the difference is greater than 31 days, aggregate by month
	if daysDiff > 31 {
		monthlyData := make(map[string]dto.UnitsSoldPerOrder)

		for _, record := range data {
			// Parse IntervalStart to time.Time for monthly aggregation
			intervalTime := record.IntervalStart

			// Format the date to "YYYY-MM" for monthly aggregation
			monthKey := intervalTime.Format("2006-01")

			if existingRecord, exists := monthlyData[monthKey]; exists {
				// Aggregate the values
				existingRecord.TotalUnitsSold += record.TotalUnitsSold
				existingRecord.TotalOrders += record.TotalOrders
				monthlyData[monthKey] = existingRecord
			} else {
				// Create a new entry in the map
				monthlyData[monthKey] = dto.UnitsSoldPerOrder{
					IntervalStart:  intervalTime, // ISO 8601 format
					TotalUnitsSold: record.TotalUnitsSold,
					TotalOrders:    record.TotalOrders,
				}
			}
		}

		// Convert the map to a slice of dto.AOVIntervalResultResponse
		result := make([]dto.UnitsSoldPerOrder, 0, len(monthlyData))
		for _, record := range monthlyData {
			// Handle NaN values by replacing with zero
			if math.IsNaN(record.TotalUnitsSold) {
				record.TotalUnitsSold = 0
			}
			if math.IsNaN(record.TotalOrders) {
				record.TotalOrders = 0
			}
			result = append(result, record)
		}
		return result, nil
	}

	// If the difference is 31 days or less, return the original data
	return data, nil
}
