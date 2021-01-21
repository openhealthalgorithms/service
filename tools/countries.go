package tools

// Country object
type Country struct {
	Name   string
	Region string
}

// CountryCode map
type CountryCode map[string]Country

// Countries returns list of countries
func Countries() CountryCode {
	countries := make(CountryCode)

	countries["AF"] = Country{"Afghanistan", "EMRD"}
	countries["AL"] = Country{"Albania", "EURB"}
	countries["DZ"] = Country{"Algeria", "AFRD"}
	countries["AS"] = Country{"Samoa", "WPRB"}
	countries["AO"] = Country{"Angola", "AFRD"}
	countries["AG"] = Country{"Antigua and Barbuda", "AMRB"}
	countries["AD"] = Country{"Andorra", "EURA"}
	countries["AR"] = Country{"Argentina", "AMRB"}
	countries["AM"] = Country{"Armenia", "EURB"}
	countries["AU"] = Country{"Australia", "#N/A"}
	countries["AT"] = Country{"Austria", "EURA"}
	countries["AZ"] = Country{"Azerbaijan", "EURB"}

	countries["BS"] = Country{"Bahamas", "AMRB"}
	countries["BH"] = Country{"Bahrain", "EMRB"}
	countries["BD"] = Country{"Bangladesh", "SEARD"}
	countries["BB"] = Country{"Barbados", "AMRB"}
	countries["BY"] = Country{"Belarus", "EURC"}
	countries["BE"] = Country{"Belgium", "EURA"}
	countries["BZ"] = Country{"Belize", "AMRB"}
	countries["BJ"] = Country{"Benin", "AFRD"}
	countries["BT"] = Country{"Bhutan", "SEARD"}
	countries["BO"] = Country{"Bolivia", "AMRD"}
	countries["BA"] = Country{"Bosnia And Herzegovina", "EURB"}
	countries["BW"] = Country{"Botswana", "AFRE"}
	countries["BR"] = Country{"Brazil", "AMRB"}
	countries["BG"] = Country{"Bulgaria", "EURB"}
	countries["BF"] = Country{"Burkina Faso", "AFRD"}
	countries["BI"] = Country{"Burundi", "AFRE"}

	countries["KH"] = Country{"Cambodia", "WPRB"}
	countries["CM"] = Country{"Cameroon", "AFRD"}
	countries["CA"] = Country{"Canada", "AMRA"}
	countries["CV"] = Country{"Cape Verde", "AFRD"}
	countries["CF"] = Country{"Central African Republic", "AFRE"}
	countries["TD"] = Country{"Chad", "AFRD"}
	countries["CL"] = Country{"Chile", "AMRB"}
	countries["CN"] = Country{"China", "WPRB"}
	countries["CO"] = Country{"Colombia", "AMRB"}
	countries["KM"] = Country{"Comoros", "AFRD"}
	countries["CD"] = Country{"Congo", "AFRE"}
	countries["CK"] = Country{"Cook Islands", "WPRB"}
	countries["CR"] = Country{"Costa Rica", "AMRB"}
	countries["CI"] = Country{"Côte d'Ivoire", "AFRE"}
	countries["HR"] = Country{"Croatia", "EURA"}
	countries["CU"] = Country{"Cuba", "AMRA"}
	countries["CY"] = Country{"Cyprus", "EURA"}
	countries["CZ"] = Country{"Czech Republic", "EURA"}

	countries["DK"] = Country{"Denmark", "EURA"}
	countries["DJ"] = Country{"Djibouti", "EMRD"}
	countries["DM"] = Country{"Dominica", "AMRB"}
	countries["DO"] = Country{"Dominican Republic", "AMRB"}

	countries["EC"] = Country{"Ecuador", "AMRD"}
	countries["EG"] = Country{"Egypt", "EMRD"}
	countries["SV"] = Country{"El Salvador", "AMRB"}
	countries["GQ"] = Country{"Equatorial Guinea", "AFRD"}
	countries["ER"] = Country{"Eritrea", "AFRE"}
	countries["EE"] = Country{"Estonia", "EURC"}
	countries["ET"] = Country{"Ethiopia", "AFRE"}

	countries["FJ"] = Country{"Fiji", "WPRB"}
	countries["FI"] = Country{"Finland", "EURA"}
	countries["FR"] = Country{"France", "EURA"}

	countries["GA"] = Country{"Gabon", "AFRD"}
	countries["GM"] = Country{"Gambia", "AFRD"}
	countries["GE"] = Country{"Georgia", "EURB"}
	countries["DE"] = Country{"Germany", "EURA"}
	countries["GH"] = Country{"Ghana", "AFRD"}
	countries["GR"] = Country{"Greece", "EURA"}
	countries["GD"] = Country{"Grenada", "AMRB"}
	countries["GT"] = Country{"Guatemala", "AMRD"}
	countries["GN"] = Country{"Guinea", "AFRD"}
	countries["GW"] = Country{"Guinea - Bissau", "AFRD"}
	countries["GY"] = Country{"Guyana", "AMRB"}

	countries["HT"] = Country{"Haiti", "AMRD"}
	countries["HN"] = Country{"Honduras", "AMRB"}
	countries["HU"] = Country{"Hungary", "EURC"}

	countries["IS"] = Country{"Iceland", "EURA"}
	countries["IN"] = Country{"India", "SEARD"}
	countries["ID"] = Country{"Indonesia", "SEARB"}
	countries["IR"] = Country{"Iran", "EMRB"}
	countries["IQ"] = Country{"Iraq", "EMRD"}
	countries["IE"] = Country{"Ireland", "EURA"}
	countries["IL"] = Country{"Israel", "EURA"}
	countries["IT"] = Country{"Italy", "EURA"}

	countries["JM"] = Country{"Jamaica", "AMRB"}
	countries["JP"] = Country{"Japan", "#N/A"}
	countries["JO"] = Country{"Jordan", "EMRB"}

	countries["KZ"] = Country{"Kazakhstan", "EURC"}
	countries["KE"] = Country{"Kenya", "AFRE"}
	countries["KI"] = Country{"Kiribati", "WPRB"}
	countries["KP"] = Country{"Korea(North)", "WPRB"}
	countries["KR"] = Country{"Korea(South)", "WPRB"}
	countries["KW"] = Country{"Kuwait", "EMRB"}
	countries["KG"] = Country{"Kyrgyzstan", "EURB"}

	countries["LA"] = Country{"Lao PDR", "WPRB"}
	countries["LV"] = Country{"Latvia", "EURC"}
	countries["LB"] = Country{"Lebanon", "EMRB"}
	countries["LS"] = Country{"Lesotho", "AFRE"}
	countries["LR"] = Country{"Liberia", "AFRD"}
	countries["LY"] = Country{"Libya", "EMRB"}
	countries["LT"] = Country{"Lithuania", "EURC"}
	countries["LU"] = Country{"Luxembourg", "EURA"}

	countries["MG"] = Country{"Madagascar", "AFRD"}
	countries["MW"] = Country{"Malawi", "AFRE"}
	countries["MY"] = Country{"Malaysia", "WPRB"}
	countries["MV"] = Country{"Maldives", "SEARD"}
	countries["ML"] = Country{"Mali", "AFRD"}
	countries["MT"] = Country{"Malta", "EURA"}
	countries["MH"] = Country{"Marshall Islands", "WPRB"}
	countries["MR"] = Country{"Mauritania", "AFRD"}
	countries["MU"] = Country{"Mauritius", "AFRD"}
	countries["MX"] = Country{"Mexico", "AMRB"}
	countries["FM"] = Country{"Micronesia (Federated States of)", "WPRB"}
	countries["MD"] = Country{"Moldova", "EURC"}
	countries["MC"] = Country{"Monaco", "EURA"}
	countries["MN"] = Country{"Mongolia", "WPRB"}
	countries["ME"] = Country{"Montenegro", "EURB"}
	countries["MA"] = Country{"Morocco", "EMRD"}
	countries["MZ"] = Country{"Mozambique", "AFRE"}
	countries["MM"] = Country{"Myanmar", "SEARD"}

	countries["NA"] = Country{"Namibia", "AFRE"}
	countries["NR"] = Country{"Nauru", "WPRB"}
	countries["NP"] = Country{"Nepal", "SEARD"}
	countries["NL"] = Country{"Netherlands", "EURA"}
	countries["NZ"] = Country{"New Zealand", "#N/A"}
	countries["NI"] = Country{"Nicaragua", "AMRD"}
	countries["NE"] = Country{"Niger", "AFRD"}
	countries["NG"] = Country{"Nigeria", "AFRD"}
	countries["NU"] = Country{"Niue", "WPRB"}
	countries["NO"] = Country{"Norway", "EURA"}

	countries["OM"] = Country{"Oman", "EMRB"}

	countries["PK"] = Country{"Pakistan", "EMRD"}
	countries["PW"] = Country{"Palau", "WPRB"}
	countries["PA"] = Country{"Panama", "AMRB"}
	countries["PG"] = Country{"Papua New Guinea", "WPRB"}
	countries["PY"] = Country{"Paraguay", "AMRB"}
	countries["PE"] = Country{"Peru", "AMRD"}
	countries["PH"] = Country{"Philippines", "WPRB"}
	countries["PL"] = Country{"Poland", "EURB"}
	countries["PT"] = Country{"Portugal", "EURA"}

	countries["QA"] = Country{"Qatar", "EMRB"}

	countries["MK"] = Country{"Macedonia, Republic of", "EURB"}
	countries["RO"] = Country{"Romania", "EURB"}
	countries["RU"] = Country{"Russian Federation", "EURC"}
	countries["RW"] = Country{"Rwanda", "AFRE"}

	countries["KN"] = Country{"Saint Kitts And Nevis", "AMRB"}
	countries["LC"] = Country{"Saint Lucia", "AMRB"}
	countries["VC"] = Country{"Saint Vincent and The Grenadines", "AMRB"}
	countries["WS"] = Country{"Samoa", "WPRB"}
	countries["SM"] = Country{"San Marino", "EURA"}
	countries["ST"] = Country{"Sao Tome and Principe", "AFRD"}
	countries["SA"] = Country{"Saudi Arabia", "EMRB"}
	countries["SN"] = Country{"Senegal", "AFRD"}
	countries["RS"] = Country{"Serbia", "EURB"}
	countries["SC"] = Country{"Seychelles", "AFRD"}
	countries["SL"] = Country{"Sierra Leone", "AFRD"}
	countries["SG"] = Country{"Singapore", "#N/A"}
	countries["SK"] = Country{"Slovakia", "EURB"}
	countries["SI"] = Country{"Slovenia", "EURA"}
	countries["SB"] = Country{"Solomon Islands", "WPRB"}
	countries["SO"] = Country{"Somalia", "EMRD"}
	countries["ZA"] = Country{"South Africa", "AFRE"}
	countries["GS"] = Country{"South Georgia and the South Sandwich Islands", "EURB"}
	countries["SS"] = Country{"South Sudan", "EMRD"}
	countries["ES"] = Country{"Spain", "EURA"}
	countries["LK"] = Country{"Sri Lanka", "SEARB"}
	countries["SD"] = Country{"Sudan", "EMRD"}
	countries["SR"] = Country{"Suriname", "AMRB"}
	countries["SZ"] = Country{"Swaziland", "AFRE"}
	countries["SE"] = Country{"Sweden", "EURA"}
	countries["CH"] = Country{"Switzerland", "EURA"}
	countries["SY"] = Country{"Syria", "EMRB"}

	countries["TW"] = Country{"Taiwan, Republic of China", "#N/A"}
	countries["TJ"] = Country{"Tajikistan", "EURB"}
	countries["TZ"] = Country{"United Republic of Tanzania", "AFRE"}
	countries["TH"] = Country{"Thailand", "SEARB"}
	countries["TG"] = Country{"Togo", "AFRD"}
	countries["TO"] = Country{"Tonga", "WPRB"}
	countries["TT"] = Country{"Trinidad and Tobago", "AMRB"}
	countries["TN"] = Country{"Tunisia", "EMRB"}
	countries["TR"] = Country{"Turkey", "EURB"}
	countries["TM"] = Country{"Turkmenistan", "EURB"}
	countries["TV"] = Country{"Tuvalu", "WPRB"}

	countries["UG"] = Country{"Uganda", "AFRE"}
	countries["UA"] = Country{"Ukraine", "EURC"}
	countries["AE"] = Country{"United Arab Emirates", "EMRB"}
	countries["GB"] = Country{"United Kingdom", "EURA"}
	countries["US"] = Country{"United States of America", "AMRA"}
	countries["UY"] = Country{"Uruguay", "AMRB"}
	countries["UZ"] = Country{"Uzbekistan", "EURB"}

	countries["VU"] = Country{"Vanuatu", "WPRB"}
	countries["VE"] = Country{"Venezuela", "AMRB"}
	countries["VN"] = Country{"Viet Nam", "WPRB"}

	countries["YE"] = Country{"Yemen", "EMRD"}

	countries["ZM"] = Country{"Zambia", "AFRE"}
	countries["ZW"] = Country{"Zimbabwe", "AFRE"}

	return countries
}
