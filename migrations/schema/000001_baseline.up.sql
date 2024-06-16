-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';


-- -----------------------------------------------------
-- Table `country`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `country` (
  `country_code` VARCHAR(3) NOT NULL,
  `country_name` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`country_code`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `currency`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `currency` (
  `currency_code` VARCHAR(3) NOT NULL,
  `currency_name` VARCHAR(50) NOT NULL,
  PRIMARY KEY (`currency_code`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `country_currency`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `country_currency` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `country_code` VARCHAR(3) NOT NULL,
  `currency_code` VARCHAR(3) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `country_code` (`country_code` ASC) VISIBLE,
  INDEX `currency_code` (`currency_code` ASC) VISIBLE,
  CONSTRAINT `country_currency_ibfk_1`
    FOREIGN KEY (`country_code`)
    REFERENCES `country` (`country_code`),
  CONSTRAINT `country_currency_ibfk_2`
    FOREIGN KEY (`currency_code`)
    REFERENCES `currency` (`currency_code`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `loan_application`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `loan_application` (
  `application_id` VARCHAR(50) NOT NULL,
  `loan_amount` DECIMAL(15,2) NOT NULL,
  `currency_code` VARCHAR(3) NOT NULL,
  `interest_rate` DECIMAL(5,2) NOT NULL,
  `loan_term` INT NOT NULL,
  `loan_term_unit` VARCHAR(10) NOT NULL,
  `application_date` DATE NOT NULL,
  `status` VARCHAR(50) NOT NULL,
  `approved_date` DATE NULL DEFAULT NULL,
  `rejection_reason` TEXT NULL DEFAULT NULL,
  `repayment_start_date` DATE NULL DEFAULT NULL,
  `approved_amount` DECIMAL(15,2) NULL,
  `income` DECIMAL(15,2) NOT NULL,
  `credit_score` INT NOT NULL,
  `existing_debts` DECIMAL(15,2) NOT NULL,
  `country_code` VARCHAR(3) NOT NULL,
  `eligible_loan_amount` DECIMAL(15,2) NULL DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`application_id`),
  INDEX `idx_currency_code` (`currency_code` ASC) VISIBLE,
  INDEX `idx_country_code` (`country_code` ASC) VISIBLE,
  CONSTRAINT `loan_application_ibfk_1`
    FOREIGN KEY (`currency_code`)
    REFERENCES `currency` (`currency_code`),
  CONSTRAINT `loan_application_ibfk_2`
    FOREIGN KEY (`country_code`)
    REFERENCES `country` (`country_code`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `user` (
  `user_id` VARCHAR(50) NOT NULL,
  `user_name` VARCHAR(255) NOT NULL,
  `user_email` VARCHAR(255) NULL,
  `user_password` VARCHAR(255) NULL,
  `user_type` VARCHAR(50) NOT NULL,
  `mobile_number` VARCHAR(15) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `user_email` (`user_email` ASC, `user_type` ASC) VISIBLE,
  UNIQUE INDEX `mobile_number` (`mobile_number` ASC, `user_type` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `loan_application_participant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `loan_application_participant` (
  `participant_id` VARCHAR(50) NOT NULL,
  `participant_type` VARCHAR(50) NOT NULL,
  `application_id` VARCHAR(50) NOT NULL,
  `user_id` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`participant_id`),
  INDEX `fk_loan_application_participant_loan_application1_idx` (`application_id` ASC) VISIBLE,
  INDEX `fk_loan_application_participant_user1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_loan_application_participant_loan_application1`
    FOREIGN KEY (`application_id`)
    REFERENCES `loan_application` (`application_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_loan_application_participant_user1`
    FOREIGN KEY (`user_id`)
    REFERENCES `user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `loan_eligibility_config`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `loan_eligibility_config` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `country_code` VARCHAR(3) NOT NULL,
  `min_credit_score` INT NOT NULL,
  `max_credit_score` INT NOT NULL,
  `max_foir` DECIMAL(5,2) NOT NULL,
  `base_loan_amount` DECIMAL(15,2) NOT NULL,
  `credit_score_factor_high` DECIMAL(5,2) NOT NULL,
  `credit_score_factor_medium` DECIMAL(5,2) NOT NULL,
  `credit_score_factor_low` DECIMAL(5,2) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_country_code` (`country_code` ASC) VISIBLE,
  CONSTRAINT `loan_eligibility_config_ibfk_1`
    FOREIGN KEY (`country_code`)
    REFERENCES `country` (`country_code`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `repayment`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `repayment` (
  `repayment_id` VARCHAR(50) NOT NULL,
  `installment_number` INT NOT NULL,
  `installment_date` DATE NOT NULL,
  `payment_date` DATE NULL DEFAULT NULL,
  `principle_amount` DECIMAL(15,2) NOT NULL,
  `interest_amount` DECIMAL(15,2) NOT NULL,
  `amount_due` DECIMAL(15,2) NOT NULL,
  `amount_paid` DECIMAL(15,2) NULL DEFAULT '0.00',
  `outstanding_balance` DECIMAL(15,2) NULL,
  `application_id` VARCHAR(50) NOT NULL,
  `status` VARCHAR(50) NOT NULL,
  `payment_reference` VARCHAR(255) NULL DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`repayment_id`),
  INDEX `idx_due_date` (`installment_date` ASC) VISIBLE,
  INDEX `idx_status` (`status` ASC) VISIBLE,
  INDEX `idx_payment_reference` (`payment_reference` ASC) VISIBLE,
  INDEX `fk_repayment_loan_application1_idx` (`application_id` ASC) VISIBLE,
  CONSTRAINT `fk_repayment_loan_application1`
    FOREIGN KEY (`application_id`)
    REFERENCES `loan_application` (`application_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `user_kyc`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_kyc` (
  `kyc_id` VARCHAR(50) NOT NULL,
  `kyc_type` VARCHAR(50) NOT NULL,
  `kyc_number` VARCHAR(50) NOT NULL,
  `user_id` VARCHAR(50) NOT NULL,
  `country_code` VARCHAR(3) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`kyc_id`),
  INDEX `idx_identification_type` (`kyc_type` ASC) VISIBLE,
  INDEX `idx_country_code` (`country_code` ASC) VISIBLE,
  INDEX `fk_user_kyc_users1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `user_identification_ibfk_2`
    FOREIGN KEY (`country_code`)
    REFERENCES `country` (`country_code`),
  CONSTRAINT `fk_user_kyc_users1`
    FOREIGN KEY (`user_id`)
    REFERENCES `user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `payment`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `payment` (
  `payment_id` VARCHAR(50) NOT NULL,
  `currency_code` VARCHAR(3) NOT NULL,
  `amount` DECIMAL(15,2) NOT NULL,
  `application_id` VARCHAR(50) NOT NULL,
  `status` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`payment_id`),
  INDEX `fk_payment_currency1_idx` (`currency_code` ASC) VISIBLE,
  INDEX `fk_payment_loan_application1_idx` (`application_id` ASC) VISIBLE,
  CONSTRAINT `fk_payment_currency1`
    FOREIGN KEY (`currency_code`)
    REFERENCES `currency` (`currency_code`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_payment_loan_application1`
    FOREIGN KEY (`application_id`)
    REFERENCES `loan_application` (`application_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS aspire_lms.repayment_payment_log (
    log_id VARCHAR(50) NOT NULL,
    repayment_id VARCHAR(50) NOT NULL,
    payment_id VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (log_id),
    INDEX fk_repayment_payment_log_repayment1_idx (repayment_id ASC) VISIBLE,
    INDEX fk_repayment_payment_log_payment1_idx (payment_id ASC) VISIBLE,
    CONSTRAINT fk_repayment_payment_log_repayment1
    FOREIGN KEY (repayment_id)
    REFERENCES aspire_lms.repayment (repayment_id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
    CONSTRAINT fk_repayment_payment_log_payment1
    FOREIGN KEY (payment_id)
    REFERENCES aspire_lms.payment (payment_id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
    ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
