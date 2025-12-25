func (r *ShipmentRepository) CalculateFreight(
    ctx context.Context,
    courierCompany string,
    weight, volume float64,
) (*FreightResult, error) {
    var result FreightResult

    err := r.db.WithContext(ctx).
        Raw(`SELECT * FROM calculate_freight(?, ?, ?)`,
            courierCompany, weight, volume).
        Scan(&result).Error

    return &result, err
}
