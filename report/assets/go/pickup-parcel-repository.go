func (r *ParcelRepository) PickupParcel(ctx context.Context, pickupCode, phone string) error {
    return r.tm.Transaction(ctx, func(ctx context.Context) error {
        // 使用行级锁查询包裹
        var parcel model.Parcel
        err := r.db.WithContext(ctx).
            Where("pickup_code = ? AND deleted_at IS NULL", pickupCode).
            Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&parcel).Error

        if err != nil {
            return ErrParcelNotFound
        }

        // 验证手机号
        if parcel.RecipientPhone != phone {
            return ErrPhoneMismatch
        }

        // 调用存储过程处理取件
        _, err = r.db.WithContext(ctx).
            Exec("SELECT * FROM pickup_parcel(?, ?)", pickupCode, phone).
            Error

        return err
    })
}
