Пример запроса постраничного списка

```go
package app

// Paginated Постраничный список
func (repo *LoyaltyBannerRepository) Paginated(page int, perPage int) (*pagination.Paginated[models.LoyaltyBanner], error) {
	helper := helpers.NewGormPaginatedHelper[models.LoyaltyBanner](repo.client).SetPerPage(perPage)
	result, err := helper.Paginated(page, func(client *gorm.DB) *gorm.DB {

		return client.Where("id = ?", 10)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

```

Примеры запросо на модификацию

```go
package app

func (repo *LoyaltyBannerRepository) Create(route *models.LoyaltyBanner) (int, error) {
	helper := helpers.NewGormModifyHelper[models.LoyaltyBanner](repo.client)
	result, err := helper.Create(route)

	if err != nil {
		return 0, err
	}

	return result.Id, nil
}

func (repo *LoyaltyBannerRepository) Update(route *models.LoyaltyBanner) error {
	helper := helpers.NewGormModifyHelper[models.LoyaltyBanner](repo.client)
	err := helper.Update(route)

	if err != nil {
		return err
	}

	return nil
}

func (repo *LoyaltyBannerRepository) Delete(route *models.LoyaltyBanner) error {
	helper := helpers.NewGormModifyHelper[models.LoyaltyBanner](repo.client)
	err := helper.Delete(route)

	if err != nil {
		return err
	}

	return nil
}

```