func (r ParcelRepository) ReceiveParcel(ctx context.Context, params *ReceiveParams) (ReceiveResult, error) {
    var result ReceiveResult

    err := r.db.WithContext(ctx).
        Raw(`SELECT * FROM receive_parcel(?, ?, ?, ?, ?, ?)`,
            params.TrackingNumber,
            params.RecipientName,
            params.RecipientPhone,
            params.CourierCompany,
            params.Size,
            params.Weight,
        ).
        Scan(&result).Error

    return &result, err
}
