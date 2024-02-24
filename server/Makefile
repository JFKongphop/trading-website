build:
	go build -o bin/main main.go

run:
	nodemon --exec go run main.go --signal SIGTERM

line:
	cloc *

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

# /stock path request
createst:
	sh ${STOCK_DIR}createst.sh

collection:
	sh ${STOCK_DIR}collection.sh

collections:
	sh ${STOCK_DIR}collections.sh

createor:
	sh ${STOCK_DIR}createor.sh

delst:
	sh ${STOCK_DIR}delst.sh

editn:
	sh ${STOCK_DIR}editn.sh

edits:
	sh ${STOCK_DIR}edits.sh

getp:
	sh ${STOCK_DIR}getp.sh

setp:
	sh ${STOCK_DIR}setp.sh

top:
	sh ${STOCK_DIR}top.sh

tran:
	sh ${STOCK_DIR}transaction.sh
