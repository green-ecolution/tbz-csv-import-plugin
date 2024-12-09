package internal

import (
  "github.com/green-ecolution/green-ecolution-backend/client"
)

type GreenEcolutionRepo interface {
  GetInfo() client.ApiGetAppInfoRequest
}

type CsvRepo interface {
}
