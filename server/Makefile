build:
	go build -o bin/main main.go

run:
	nodemon --exec go run main.go --signal SIGTERM

USER_DIR := ./bash-request/user/
STOCK_DIR := ./bash-request/stock/

# /user path request
signup:
	sh ${USER_DIR}signup.sh

deposit:
	sh ${USER_DIR}deposit.sh

withdraw:
	sh ${USER_DIR}withdraw.sh

buy:
	sh ${USER_DIR}buy.sh

sale:
	sh ${USER_DIR}sale.sh

setfav:
	sh ${USER_DIR}setfav.sh

balance:
	sh ${USER_DIR}balance.sh

balancehis:
	sh ${USER_DIR}balancehis.sh

getfav:
	sh ${USER_DIR}getfav.sh

signin:
	sh ${USER_DIR}signin.sh

history:
	sh ${USER_DIR}history.sh

ustock:
	sh ${USER_DIR}ustock.sh

ustockamount:
	sh ${USER_DIR}ustockamount.sh

delfav:
	sh ${USER_DIR}delfav.sh

delacc:
	sh ${USER_DIR}delacc.sh