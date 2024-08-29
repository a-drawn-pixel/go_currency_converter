package Application

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type CachingService struct {
	next  CurrencyService
	cache *cache.Cache
}

func NewCachingService(next CurrencyService, cacheExpiry time.Duration) *CachingService {
	c := cache.New(cacheExpiry, cacheExpiry)
	return &CachingService{
		next:  next,
		cache: c,
	}
}

func (s *CachingService) GetRates(currency string) (map[string]float64, error) {
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	if cachedRates, exists := s.cache.Get(currency); exists {
		fmt.Println("cached rate for", currency)
		return cachedRates.(map[string]float64), nil
	}

	rates, err := s.next.GetRates(currency)
	if err != nil {
		return nil, err
	}

	s.cache.Set(currency, rates, cache.DefaultExpiration)
	fmt.Println("caching rate for", currency)
	return rates, nil
}
