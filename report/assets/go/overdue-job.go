func (j *OverdueJob) Process(ctx context.Context) error {
    // 查询滞留件视图
    var overdues []model.OverdueParcel
    err := j.db.WithContext(ctx).
        Table("v_pending_overdue_parcels").
        Where("overdue_days > 0").
        Find(&overdues).Error

    if err != nil {
        return err
    }

    // 批量处理滞留件
    for _, overdue := range overdues {
        j.processOverdueParcel(ctx, &overdue)
    }

    return nil
}
