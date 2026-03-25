BEGIN;

CREATE TABLE IF NOT EXISTS currencies (
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(3) UNIQUE NOT NULL,
    name        VARCHAR(50) NOT NULL,
    symbol      VARCHAR(10) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_currencies_code ON currencies(code);

INSERT INTO currencies (code, name, symbol, is_active)
VALUES
('USD','US Dollar','$',true),
('EUR','Euro','€',true),
('JPY','Japanese Yen','¥',true),
('GBP','British Pound Sterling','£',true),
('CHF','Swiss Franc','CHF',true),
('AUD','Australian Dollar','$',true),
('CAD','Canadian Dollar','$',true),
('NZD','New Zealand Dollar','$',true),
('CNY','Chinese Yuan','¥',true),
('HKD','Hong Kong Dollar','$',true),
('SGD','Singapore Dollar','$',true),
('INR','Indian Rupee','₹',true),
('KRW','South Korean Won','₩',true),
('THB','Thai Baht','฿',true),
('MYR','Malaysian Ringgit','RM',true),
('IDR','Indonesian Rupiah','Rp',true),
('PHP','Philippine Peso','₱',true),
('VND','Vietnamese Dong','₫',true),
('PKR','Pakistani Rupee','₨',true),
('BDT','Bangladeshi Taka','৳',true),
('AED','UAE Dirham','د.إ',true),
('SAR','Saudi Riyal','﷼',true),
('QAR','Qatari Riyal','﷼',true),
('KWD','Kuwaiti Dinar','د.ك',true),
('BHD','Bahraini Dinar','ب.د',true),
('OMR','Omani Rial','﷼',true),
('ILS','Israeli Shekel','₪',true),
('TRY','Turkish Lira','₺',true),
('SEK','Swedish Krona','kr',true),
('NOK','Norwegian Krone','kr',true),
('DKK','Danish Krone','kr',true),
('PLN','Polish Zloty','zł',true),
('CZK','Czech Koruna','Kč',true),
('HUF','Hungarian Forint','Ft',true),
('RON','Romanian Leu','lei',true),
('BGN','Bulgarian Lev','лв',true),
('HRK','Croatian Kuna','kn',true),
('BRL','Brazilian Real','R$',true),
('MXN','Mexican Peso','$',true),
('ARS','Argentine Peso','$',true),
('CLP','Chilean Peso','$',true),
('COP','Colombian Peso','$',true),
('PEN','Peruvian Sol','S/',true),
('ZAR','South African Rand','R',true),
('NGN','Nigerian Naira','₦',true),
('EGP','Egyptian Pound','£',true),
('KES','Kenyan Shilling','KSh',true),
('GHS','Ghanaian Cedi','₵',true),
('MAD','Moroccan Dirham','د.م.',true);

COMMIT;