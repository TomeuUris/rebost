package main

// ToProductDTO converts a Product domain model to a ProductDTO.
func ToProductDTO(product Product) ProductDTO {
	return ProductDTO{
		Barcode:     product.Barcode,
		ProductName: product.ProductName,
		Brand:       product.Brand,
		Ingredients: product.Ingredients,
		Nutriments:  ToNutrimentsDTO(product.Nutriments),
	}
}

// ToNutrimentsDTO converts a Nutriments domain model to a NutrimentsDTO.
func ToNutrimentsDTO(nutriments Nutriments) NutrimentsDTO {
	return NutrimentsDTO{
		Per100g:    ToRationDTO(nutriments.Per100g),
		PerServing: ToRationDTO(nutriments.PerServing),
	}
}

// ToRationDTO converts a Ration domain model to a RationDTO.
func ToRationDTO(ration *Ration) *RationDTO {
	if ration == nil {
		return nil
	}
	return &RationDTO{
		EnergyKcal:    ration.EnergyKcal,
		Fat:           ration.Fat,
		SaturatedFat:  ration.SaturatedFat,
		Carbohydrates: ration.Carbohydrates,
		Sugars:        ration.Sugars,
		Fiber:         ration.Fiber,
		Proteins:      ration.Proteins,
		Salt:          ration.Salt,
		Sodium:        ration.Sodium,
	}
}

// ToInventoryItemDTO converts an InventoryItem domain model to an InventoryItemDTO.
func ToInventoryItemDTO(item InventoryItem) InventoryItemDTO {
	return InventoryItemDTO{
		ID:             item.ID,
		Product:        ToProductDTO(item.Product),
		ExpirationDate: item.ExpirationDate,
		Quantity:       item.Quantity,
		CreatedAt:      item.CreatedAt,
	}
}

// ToProductFromAPIProduct converts an APIProduct to a Product domain model.
func ToProductFromAPIProduct(apiProduct APIProduct) Product {
	return Product{
		Barcode:     apiProduct.Barcode,
		ProductName: apiProduct.ProductName,
		Brand:       apiProduct.Brand,
		Ingredients: apiProduct.Ingredients,
		Nutriments:  apiProduct.Nutriments,
	}
}

// ToProductDTOFromAPIProduct converts an APIProduct to a ProductDTO.
func ToProductDTOFromAPIProduct(apiProduct APIProduct) ProductDTO {
	return ProductDTO{
		Barcode:     apiProduct.Barcode,
		ProductName: apiProduct.ProductName,
		Brand:       apiProduct.Brand,
		Ingredients: apiProduct.Ingredients,
		Nutriments:  ToNutrimentsDTO(apiProduct.Nutriments),
	}
}
