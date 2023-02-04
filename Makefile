migrate:
	migrate -path migrations/ -database 'sqlite3://$(HOME)/.config/gonotes/gonotes.db' up
migrate-down:
	migrate -path migrations/ -database 'sqlite3://$(HOME)/.config/gonotes/gonotes.db' down
migrate-new:
	migrate create -ext sql -dir migrations -seq 'new'
gqlgen:
	go run github.com/Khan/genqlient
