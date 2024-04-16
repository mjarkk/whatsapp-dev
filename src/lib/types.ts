export interface DBModel {
	ID: number
	CreatedAt: null | string // "2024-04-18T12:19:31.792423+02:00"
	UpdatedAt: null | string // "2024-04-18T12:19:31.792423+02:00"
	DeletedAt: null | string // "2024-04-18T12:19:31.792423+02:00"
}

export const emptyDBModel = (): DBModel => ({
	ID: 0,
	CreatedAt: null,
	UpdatedAt: null,
	DeletedAt: null,
})
