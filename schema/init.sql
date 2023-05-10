CREATE TABLE users (
  userID SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(64) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(255) NOT NULL
);

CREATE TABLE contacts (
  contactID SERIAL PRIMARY KEY,
  userID INT,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL
);

CREATE TABLE accounts (
  accountID SERIAL PRIMARY KEY,
  userID INT,
  name VARCHAR(255) NOT NULL,
  balance REAL NOT NULL,
  currency VARCHAR(10) NOT NULL
);

CREATE TABLE transactions (
  transactionID SERIAL PRIMARY KEY,
  fromAccountID INT NOT NULL,
  toAccountID INT,
  date timestamp NOT NULL,
  amount REAL NOT NULL,
  description VARCHAR(255),
  categoryID INT,
  tagID INT,
  transactionTypeID INT
);

CREATE TABLE transactionTypes (
  transactionTypeID SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE Loans (
  loanID SERIAL PRIMARY KEY,
  accountID INT,
  contactID INT,
  amount REAL NOT NULL,
  date timestamp NOT NULL,
  repaymentDate timestamp NOT NULL,
  type VARCHAR(255) NOT NULL,
  status VARCHAR(255) NOT NULL
);

CREATE TABLE loanPayments (
  paymentID SERIAL PRIMARY KEY,
  loanID INT,
  date timestamp NOT NULL,
  amount REAL NOT NULL
);

CREATE TABLE userSettings (
  userID INT,
  currency VARCHAR(10) NOT NULL,
  dateRange VARCHAR(255) NOT NULL,
  notificationPreference VARCHAR(255) NOT NULL
);
CREATE TABLE categories (
  categoryID SERIAL PRIMARY KEY,
  userID INT,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE tags (
  tagID SERIAL PRIMARY KEY,
  userID INT,
  name VARCHAR(255) NOT NULL,
  icon VARCHAR(255) NOT NULL
);

CREATE TABLE notifications (
  notificationID SERIAL PRIMARY KEY,
  userID INT,
  title VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  date timestamp NOT NULL,
  isRead BOOLEAN NOT NULL
);

CREATE TABLE recurringTransactions (
  recurringTransactionID SERIAL PRIMARY KEY,
  userID INT,
  accountID INT,
  startDate timestamp NOT NULL,
  endDate timestamp,
  amount REAL NOT NULL,
  type VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  categoryID INT,
  tagID INT,
  transactioNtYpeID INT,
  frequency varCHAR(255) NOT NULL,
  frequencyINTErval INT NOT NULL
);

CREATE TABLE goals (
  goalID SERIAL PRIMARY KEY,
  userID INT,
  name VARCHAR(255) NOT NULL,
  amount REAL NOT NULL,
  startDate timestamp NOT NULL,
  endDate timestamp,
  description VARCHAR(255)
);

ALTER TABLE contacts ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE accounts ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE transactions ADD FOREIGN KEY (fromAccountID) REFERENCES accounts (accountID);

ALTER TABLE transactions ADD FOREIGN KEY (toAccountID) REFERENCES accounts (accountID);

ALTER TABLE transactions ADD FOREIGN KEY (categoryID) REFERENCES categories (categoryID);

ALTER TABLE transactions ADD FOREIGN KEY (tagID) REFERENCES tags (tagID);

ALTER TABLE recurringTransactions ADD FOREIGN KEY (tagID) REFERENCES tags (tagID);

ALTER TABLE transactions ADD FOREIGN KEY (transactionTypeID) REFERENCES transactionTypes (transactionTypeID);

ALTER TABLE loans ADD FOREIGN KEY (accountID) REFERENCES accounts (accountID);

ALTER TABLE loans ADD FOREIGN KEY (contactID) REFERENCES contacts (contactID);

ALTER TABLE loanPayments ADD FOREIGN KEY (loanID) REFERENCES loans (loanID);

ALTER TABLE userSettings ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE categories ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE tags ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE notifications ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE recurringTransactions ADD FOREIGN KEY (userID) REFERENCES users (userID);

ALTER TABLE recurringTransactions ADD FOREIGN KEY (accountID) REFERENCES accounts (accountID);

ALTER TABLE recurringTransactions ADD FOREIGN KEY (categoryID) REFERENCES categories (categoryID);

ALTER TABLE recurringTransactions ADD FOREIGN KEY (transactionTypeID) REFERENCES transactionTypes (transactionTypeID);

ALTER TABLE Goals ADD FOREIGN KEY (userID) REFERENCES users (userID);

-- INSERT INTO users (name, email, password, phone) VALUES
--   ('Ionel', 'ionel@mail.com', '1111', '012345');

select * from transactions inner join transaction_types on transactions.transaction_type_id = transaction_types.id;

-- Get the account balance
SELECT 
	COALESCE(SUM(CASE WHEN transaction_type_id = ? THEN amount ELSE -amount END), 0) 
	+ COALESCE((SELECT SUM(CASE WHEN to_account_id = ? THEN amount ELSE -amount END) 
				FROM transactions 
				WHERE deleted_at is null 
				AND (to_account_id = ? OR from_account_id = ?) 
				AND transaction_type_id = ?), 0) 
FROM transactions 
WHERE deleted_at is null 
AND to_account_id = ? 
AND transaction_type_id <> ?;
