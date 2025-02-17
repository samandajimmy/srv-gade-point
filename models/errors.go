package models

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalServerError to store internal server error message
	ErrInternalServerError = errors.New("Internal Server Error ")

	// ErrNotFound to store not found error message
	ErrNotFound = errors.New("Item tidak ditemukan")

	// ErrConflict to store conflicted error message
	ErrConflict = errors.New("Item sudah ada")

	// ErrSetVar to store setting variable error message
	ErrSetVar = errors.New("Setting variable error")

	// ErrMigrateNoChange to store migration no change error message
	ErrMigrateNoChange = errors.New("no change")

	// ErrEnvFileNF to store error env file not found error message
	ErrEnvFileNF = errors.New("env file not found! expected we have set system variable")

	// ErrBadParamInput to store bad parameter error message
	ErrBadParamInput = errors.New("Parameter yang diberikan tidak valid")

	// ErrCampaignFailed to store create campaign failed error message
	ErrCampaignFailed = errors.New("Gagal untuk membuat Promo")

	// ErrCampaignUpdateFailed to store update campaign failed error message
	ErrCampaignUpdateFailed = errors.New("Gagal untuk mengupdate Promo")

	// ErrNoCampaign to store campaign not found error message
	ErrNoCampaign = errors.New("Promo tidak tersedia")

	// ErrGetCampaign to get campaign error message
	ErrGetCampaign = errors.New("Terjadi kesalahan dalam mengambil Promo")

	// ErrGetCampaignCounter to get campaign counter error message
	ErrGetCampaignCounter = errors.New("Campaign tidak tersedia")

	// ErrCalculateFormulaCampaign to get campaign counter error message
	ErrCalculateFormulaCampaign = errors.New("Terjadi kesalahan dalam menghitung formula Promo")

	// ErrStoreCampaignTrx to get campaign counter error message
	ErrStoreCampaignTrx = errors.New("Terjadi kesalahan dalam menyimpan transaksi Promo")

	// ErrGetUserPoint to get user point error message
	ErrGetUserPoint = errors.New("Terjadi kesalahan dalam mengambil Point User")

	// ErrUserPointNA to get user point N/A error message
	ErrUserPointNA = errors.New("Anda belum memiliki Point")

	// ErrGetUserPointHistory to get user point history error message
	ErrGetUserPointHistory = errors.New("Terjadi kesalahan dalam mengambil History Point User")

	// ErrUserPointHistoryNA to get user point history N/A error message
	ErrUserPointHistoryNA = errors.New("Anda belum memiliki history point")

	// ErrCampaignExpired to store campaign expired error message
	ErrCampaignExpired = errors.New("Promo telah habis masa berlaku")

	// ErrPointDeficit to store point deficit error message
	ErrPointDeficit = errors.New("Anda belum memiliki cukup point untuk  membeli Voucher ini")

	// ErrVoucherExpired to store voucher expired error message
	ErrVoucherExpired = errors.New("Voucher telah expire")

	// ErrVoucherFailed to store create voucher failed error message
	ErrVoucherFailed = errors.New("Gagal untuk membuat voucher")

	// ErrVoucherNotFound to inquiry list voucher error message
	ErrVoucherNotFound = errors.New("Voucher tidak ditemukan")

	// ErrVoucherGenearatePromoCodes to store generate promo codes error message
	ErrVoucherGenearatePromoCodes = errors.New("Terjadi kesalahan dalam membuat Kode Promo")

	// ErrVoucherStorePomoCodes to store generate promo codes error message
	ErrVoucherStorePomoCodes = errors.New("Terjadi kesalahan dalam menyimpan kode promosi`")

	// ErrDeleteVoucher to store delete voucher error message
	ErrDeleteVoucher = errors.New("Terjadi kesalahan dalam menghapus voucher")

	// ErrVoucherUpdateFailed to store update voucher failed error message
	ErrVoucherUpdateFailed = errors.New("Gagal mengupdate Voucher")

	// ErrOpenVoucherImg to store open file image
	ErrOpenVoucherImg = errors.New("Gagal mengupload voucher")

	// ErrWriteVoucherImg to store open file image
	ErrWriteVoucherImg = errors.New("Gagal mengupload voucher")

	// ErrCopyVoucherImg to store open file image
	ErrCopyVoucherImg = errors.New("Tidak bisa mengupload file image voucher ")

	// ErrGetVouchers to get vouchers error message
	ErrGetVouchers = errors.New("Terjadi kesalahan dalam mengambil voucher")

	// ErrGetVouchers to get vouchers error message
	ErrGetRewardPromotions = errors.New("Terjadi kesalahan dalam mengambil reward promotions")

	// ErrGetVoucherCodes to get vouchers error message
	ErrGetVoucherCodes = errors.New("Terjadi kesalahan dalam mengambil kode voucher")

	// ErrGetVoucherCounter to get voucher counter error message
	ErrGetVoucherCounter = errors.New("Terjadi kesalahan dalam mengambil kode perhitungan kode voucher")

	// ErrRedeemVoucher to store redeem voucher error message
	ErrRedeemVoucher = errors.New("Kode voucher tidak tersedia untuk di redeem")

	// ErrExceedBuyLimit to store exceeded buying limit voucher error message
	ErrExceedBuyLimit = errors.New("Pembelian voucher sudah melebihi limit")

	// ErrVoucherNotStarted to store voucher not started error message
	ErrVoucherNotStarted = errors.New("Penggunaan Voucher belum dimulai")

	// ErrVoucherUnavailable to store voucher unavailable error message
	ErrVoucherUnavailable = errors.New("Voucher tidak tersedia")

	// ErrBuyingVoucherExceeded to store exceeded error message
	ErrBuyingVoucherExceeded = errors.New("Voucher tidak dapat dibeli lagi untuk hari ini")

	// ErrVoucherOutOfStock to store voucher unavailable error message
	ErrVoucherOutOfStock = errors.New("Voucher sudah habis")

	// ErrUpdatePromoCodes to update promo codes error message
	ErrUpdatePromoCodes = errors.New("Terjadi kesalahan dalam mengupdate kode promo voucher")

	// ErrVoucherPoint to get voucher point error message
	ErrVoucherPoint = errors.New("Terjadi kesalahan dalam mengambil point voucher")

	// ErrVoucherCodeUnavailable to store voucher unavailable error message
	ErrVoucherCodeUnavailable = errors.New("Kode voucher tidak tersedia")

	// ErrValidatorUnavailable to store validator unavailable error message
	ErrValidatorUnavailable = errors.New("Validator tidak tersedia")

	// ErrValidation to store validation error message
	ErrValidation = errors.New("Ada kesalahan dalam input Anda")

	// ErrValidationTrxAmt to store validation error message
	ErrValidationTrxAmt = errors.New("Jumlah transaksi Anda tidak mencukupi untuk voucher ini")

	// ErrUsername to store username error message
	ErrUsername = errors.New("Username atau Password yang digunakan tidak valid")

	// ErrPassword to store password error message
	ErrPassword = errors.New("Username atau Password yang digunakan tidak valid")

	// ErrTokenExpired to store password error message
	ErrTokenExpired = errors.New("Token Anda telah expire")

	// ErrUsersNA to store users not available error message
	ErrUsersNA = errors.New("User tidak tersedia")

	// ErrGetUsersPoint to store get users point error message
	ErrGetUsersPoint = errors.New("Terjadi kesalahan dalam mengambil point user")

	// ErrGetVoucherHistory to store get voucher history error message
	ErrGetVoucherHistory = errors.New("Terjadi kesalahan dalam mengambil history voucher")

	// ErrStartDateFormat to store a date format params error message
	ErrStartDateFormat = errors.New("Parameter tanggal tidak sesuai format")

	// ErrEndDateFormat to store a date format params error message
	ErrEndDateFormat = errors.New("Parameter tanggal tidak sesuai format")

	// ErrAllowedExtVchrCodesImport to store allowed file ext error message
	ErrAllowedExtVchrCodesImport = errors.New("Hanya menerima format file CSV dan JSON")

	// ErrMappingVchrCodesImport to store allowed file ext error message
	ErrMappingVchrCodesImport = errors.New("Terjadi kesalahan dalam mengimport mapping data dari file yang diupload")

	// ErrGetReward to store get reward error message
	ErrGetReward = errors.New("Terjadi kesalahan dalam mengambil reward")

	// ErrGetRewardCounter to get reward counter error message
	ErrGetRewardCounter = errors.New("Terjadi kesalahan dalam mengambil total reward")

	// ErrSameCifReferrerAndReferral to get reward referral error message
	ErrSameCifReferrerAndReferral = errors.New("Terjadi kesalahan cif referral dan referrer sama")

	// ErrValidateGetReferral to get reward referral error message
	ErrValidateGetReferral = errors.New("Kode referral sudah pernah digunakan")

	// ErrValidateGetReferralMaxReward to get max reward referral error message
	ErrValidateGetReferralMaxReward = errors.New("Kode referral telah melampui batas milestone rewards")

	// ErrRewardFailed to store create reward failed error message
	ErrRewardFailed = errors.New("Gagal untuk membuat reward")

	// ErrTagFailed to store create tag failed error message
	ErrTagFailed = errors.New("Gagal untuk membuat Tag")

	// ErrQuotaFailed to store create quota failed error message
	ErrQuotaFailed = errors.New("Gagal untuk membuat Kuota")

	// ErrCheckQuotaFailed to check quota failed error message
	ErrCheckQuotaFailed = errors.New("Gagal untuk mengecek Kuota")

	// ErrQuotaNotAvailable to info quota is not available message
	ErrQuotaNotAvailable = errors.New("Kuota tidak tersedia")

	// ErrQuotaNACIF to info CIF quota is not available message
	ErrQuotaNACIF = errors.New("Kuota tidak tersedia untuk CIF ini")

	// ErrTodaysQuotaNotAvailable to info quota is not available message
	ErrTodaysQuotaNotAvailable = errors.New("Kuota tidak tersedia untuk hari ini")

	// ErrCreateRewardsFailed to store create rewards failed message
	ErrCreateRewardsFailed = errors.New("Terjadi kesalahan dalam membuat Reward")

	// ErrRewardTrxFailed to store create reward transaction failed error message
	ErrRewardTrxFailed = errors.New("Gagal dalam membuat transaksi reward")

	// ErrRewardTrxUpdateFailed to store create reward transaction failed error message
	ErrRewardTrxUpdateFailed = errors.New("Gagal mengupdate transaksi reward")

	// ErrDelRewardFailed to store delete reward error message
	ErrDelRewardFailed = errors.New("Terjadi kesalahan dalam menghapus reward")

	// ErrDelQuotaFailed to store delete quota error message
	ErrDelQuotaFailed = errors.New("Terjadi kesalahan dalam mengupdate Kuota")

	// ErrAddQuotaFailed to store add quota error message
	ErrAddQuotaFailed = errors.New("Terjadi kesalahan dalam menambah Kuota")

	// ErrReduceQuotaFailed to store minus quota error message
	ErrReduceQuotaFailed = errors.New("Terjadi keselahan dalam mengurangi Kuota")

	// ErrCreateQuotasFailed to store create quotas failed message
	ErrCreateQuotasFailed = errors.New("Terjadi kesalahan dalam membuat Kuota")

	// ErrRefreshQuotaFailed to store minus quota error message
	ErrRefreshQuotaFailed = errors.New("Terjadi kesalahan dalam me-refresh Kuota")

	// ErrCreateTagsFailed to store create tags failed message
	ErrCreateTagsFailed = errors.New("Terjadi kesalahan dalam membuat Tag")

	// ErrPromoCode to store promo code error message
	ErrPromoCode = errors.New("Promo Code tidak tersedia")

	// ErrTrxDateFormat to store a trx date format params error message
	ErrTrxDateFormat = errors.New("Parameter tanggal tidak sesuai dengan Format")

	// ErrRefTrxNotFound to not found ref_trx error message
	ErrRefTrxNotFound = errors.New("ID transaksi tidak ditemukan")

	// ErrGetRewardTrxCounter to get reward transaction counter error message
	ErrGetRewardTrxCounter = errors.New("Terjadi kesalahan dalam mengambil total reward transaction")

	// ErrGetRewardTrx to store get reward transaction error message
	ErrGetRewardTrx = errors.New("Terjadi kesalahan dalam mengambil reward transaction")

	// ErrMessageNoRewards to store a no rewards message response of data
	ErrMessageNoRewards = errors.New("Maaf, tidak ada reward yang tersedia")

	// ErrMessageRewardTrxAlreadyExists to store reward transaction already exists
	ErrMessageRewardTrxAlreadyExists = errors.New("Transaksi reward sudah ada")

	// ErrRefIDStatus to not found ref_trx error message
	ErrRefIDStatus = errors.New("Transaksi ID ")

	// ErrMilestone to not found ref_trx error message
	ErrMilestone = errors.New("Data milestone tidak ditemukan")

	// ErrRanking to not found ref_trx error message
	ErrRanking = errors.New("Data ranking tidak ditemukan")

	// ErrCIF to not found data cif
	ErrCIF = errors.New("CIF tidak ditemukan")

	// ErrVoucherPrcessed voucher still on progress
	ErrVoucherPrcessed = errors.New("Kode voucher sedang diproses")

	// ErrVoucherRedemeed voucher voucher has been redemmeed
	ErrVoucherRedemeed = errors.New("Kode voucher sudah digunakan")

	ErrFormulaNotValid = errors.New("formula is not valid")

	// ErrCoreFailed to store create core Transaction failed error message
	ErrCoreTrxFailed = errors.New("Gagal untuk menyimpan transaksi cicil emas schedule")

	// ErrCoreTrxEmpty to store create core Transaction empty data
	ErrCoreTrxEmpty = errors.New("Tidak ada data untuk di simpan")

	// ErrGenerateRefCodeFailed to store create reward failed error message
	ErrGenerateRefCodeFailed = errors.New("Gagal untuk membuat referral code")

	// ErrRefCampaginNF to store create reward failed error message
	ErrRefCampaginNF = errors.New("Campaign referral tidak tersedia untuk")

	// ErrRefCodesNF to store get referral codes error
	ErrRefCodesNF = errors.New("Kode referral tidak ditemukan")

	// ErrRefTrxExceeded to store exceeded referral code trx error message
	ErrRefTrxExceeded = errors.New("Kode referral tidak dapat digunakan")

	// ErrRefTrxIncentive to store max incentive days referral code trx error message
	ErrRefTrxIncentive = errors.New("Maaf, kode referral telah mencapai maksimal insentif")

	// ErrRefTrxExceededPC to store exceeded referral code trx error message when used on same product code
	ErrRefTrxExceededPC = errors.New("Kode referral tidak dapat digunakan pada produk yang sama")

	// ErrXpoinApi to store xpoin api error message
	ErrXpoinApi = errors.New("Terjadi kesalahan pada Xpoin")

	// ErrIncDayExceeded to store exceeded incentive per day error message
	DynErrIncDayExceeded = "Saldo insentif per hari maksimal Rp %s"

	// ErrIncMonthExceeded to store exceeded incentive per day error message
	DynErrIncMonthExceeded = "Saldo insentif per bulan maksimal Rp %s"

	// ErrFormulaSetup to store formula setup error message
	ErrFormulaSetup = errors.New("Formula tidak lengkap")

	// ErrRefPrefixNF to store get referral codes error
	ErrRefPrefixNF = errors.New("Prefix untuk promo code tersebut tidak ditemukan")

	// ErrRefHistoryIncentiveNF to store get history incentive error message
	ErrRefHistoryIncentiveNF = errors.New("Data history incentive referral tidak ditemukan")

	// ErrFriends get friends by cif data not found
	ErrFriends = errors.New("Data Sahabat tidak di temukan")

	// ErrOslInactive error message of get osl status
	ErrOslInactive = errors.New("OSL dari nasabah telah aktif")

	// ErrXpoinAPIRequest to store xpoin api request error message
	ErrXpoinAPIRequest = "XPOIN API: RC-%s - %s"

	// ErrDynRewardValidation to store dynamic error message of rewards validation
	ErrDynRewardValidation = "RewardIdx #%s: %s"
)

// DynamicErr to return parameterize errors
func DynamicErr(message string, args ...interface{}) error {

	return fmt.Errorf(message, args...)
}
