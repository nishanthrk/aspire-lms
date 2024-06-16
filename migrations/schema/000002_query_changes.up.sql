
-- -----------------------------------------------------
-- Data `currency`
-- -----------------------------------------------------
INSERT INTO currency (currency_code, currency_name) VALUES
('USD', 'US Dollar'),
('INR', 'Indian Rupee'),
('CAD', 'Canadian Dollar'),
('GBP', 'British Pound Sterling'),
('AUD', 'Australian Dollar'),
('EUR', 'Euro'),
('JPY', 'Japanese Yen'),
('CNY', 'Chinese Yuan'),
('BRL', 'Brazilian Real'),
('RUB', 'Russian Ruble'),
('ZAR', 'South African Rand');

-- -----------------------------------------------------
-- Data `country`
-- -----------------------------------------------------
INSERT INTO country (country_code, country_name) VALUES
('USA', 'United States of America'),
('IND', 'India'),
('CAN', 'Canada'),
('GBR', 'United Kingdom'),
('AUS', 'Australia'),
('DEU', 'Germany'),
('FRA', 'France'),
('JPN', 'Japan'),
('CHN', 'China'),
('BRA', 'Brazil'),
('RUS', 'Russia'),
('ZAF', 'South Africa');

-- -----------------------------------------------------
-- Data `country_currency`
-- -----------------------------------------------------
INSERT INTO country_currency (country_code, currency_code) VALUES
('USA', 'USD'),
('IND', 'INR'),
('CAN', 'CAD'),
('GBR', 'GBP'),
('AUS', 'AUD'),
('DEU', 'EUR'),
('FRA', 'EUR'),
('JPN', 'JPY'),
('CHN', 'CNY'),
('BRA', 'BRL'),
('RUS', 'RUB'),
('ZAF', 'ZAR');

INSERT INTO `loan_eligibility_config` (`country_code`, `min_credit_score`, `max_credit_score`, `max_foir`, `base_loan_amount`, `credit_score_factor_high`, `credit_score_factor_medium`, `credit_score_factor_low`) VALUES
('USA', 650, 850, 0.5, 1000000, 1.2, 1.1, 1.0),
('IND', 600, 800, 0.45, 750000, 1.1, 1.0, 0.9),
('CAN', 620, 820, 0.48, 800000, 1.15, 1.05, 0.95),
('GBR', 640, 840, 0.46, 900000, 1.18, 1.08, 0.98),
('AUS', 630, 830, 0.47, 850000, 1.17, 1.07, 0.97),
('DEU', 650, 850, 0.44, 950000, 1.2, 1.1, 1.0),
('FRA', 620, 820, 0.49, 800000, 1.15, 1.05, 0.95),
('JPN', 600, 800, 0.45, 750000, 1.1, 1.0, 0.9),
('CHN', 630, 830, 0.47, 850000, 1.17, 1.07, 0.97),
('BRA', 620, 820, 0.48, 800000, 1.15, 1.05, 0.95),
('RUS', 650, 850, 0.5, 1000000, 1.2, 1.1, 1.0),
('ZAF', 610, 810, 0.46, 780000, 1.13, 1.03, 0.93);

INSERT INTO `user` (`user_id`, `user_name`, `user_email`, `user_password`, `user_type`, `mobile_number`, `created_at`, `updated_at`)
VALUES ('53297921-01d9-4311-94f3-54cbb971c5a0', 'Approver One', NULL, 'fa585d89c851dd338a70dcf535aa2a92fee7836dd6aff1226583e88e0996293f16bc009c652826e0fc5c706695a03cddce372f139eff4d13959da6f1f5d3eabe', 'EMPLOYEE', '9790970381', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

