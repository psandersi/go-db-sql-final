package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotNil(t, id)

	parcels, err := store.Get(id)
	require.NoError(t, err)
	assert.NotEmpty(t, parcels.Number)
	assert.Equal(t, parcel.Client, parcels.Client)
	assert.Equal(t, parcel.CreatedAt, parcels.CreatedAt)
	assert.Equal(t, parcel.Status, parcels.Status)
	assert.Equal(t, parcel.Address, parcels.Address)

	err = store.Delete(id)
	require.NoError(t, err)
	m, _ := store.Get(id)
	require.Equal(t, "", m.Address)

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotNil(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	parcels, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, parcels.Address)
	assert.NotEqual(t, parcel.Address, parcels.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotNil(t, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	parcels, err := store.Get(id)
	require.NoError(t, err)
	assert.NotEqual(t, parcel.Status, parcels.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotNil(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)

	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		assert.Equal(t, parcel.Address, parcelMap[parcel.Number].Address)
		assert.Equal(t, parcel.Client, parcelMap[parcel.Number].Client)
		assert.Equal(t, parcel.CreatedAt, parcelMap[parcel.Number].CreatedAt)
		assert.Equal(t, parcel.Status, parcelMap[parcel.Number].Status)
	}
}
