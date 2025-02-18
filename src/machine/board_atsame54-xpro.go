// +build atsame54_xpro

package machine

import (
	"device/sam"
	"runtime/interrupt"
)

// Definition for compatibility, but not used
const RESET_MAGIC_VALUE = 0x00000000

const (
	LED    = PC18
	BUTTON = PB31
)

const (
	// https://ww1.microchip.com/downloads/en/DeviceDoc/70005321A.pdf

	// Extension Header EXT1
	EXT1_PIN3_ADC_P     = PB04
	EXT1_PIN4_ADC_N     = PB05
	EXT1_PIN5_GPIO1     = PA06
	EXT1_PIN6_GPIO2     = PA07
	EXT1_PIN7_PWM_P     = PB08
	EXT1_PIN8_PWM_N     = PB09
	EXT1_PIN9_IRQ       = PB07
	EXT1_PIN9_GPIO      = PB07
	EXT1_PIN10_SPI_SS_B = PA27
	EXT1_PIN10_GPIO     = PA27
	EXT1_PIN11_TWI_SDA  = PA22
	EXT1_PIN12_TWI_SCL  = PA23
	EXT1_PIN13_UART_RX  = PA05
	EXT1_PIN14_UART_TX  = PA04
	EXT1_PIN15_SPI_SS_A = PB28
	EXT1_PIN16_SPI_SDO  = PB27
	EXT1_PIN17_SPI_SDI  = PB29
	EXT1_PIN18_SPI_SCK  = PB26

	// Extension Header EXT2
	EXT2_PIN3_ADC_P     = PB00
	EXT2_PIN4_ADC_N     = PA03
	EXT2_PIN5_GPIO1     = PB01
	EXT2_PIN6_GPIO2     = PB06
	EXT2_PIN7_PWM_P     = PB14
	EXT2_PIN8_PWM_N     = PB15
	EXT2_PIN9_IRQ       = PD00
	EXT2_PIN9_GPIO      = PD00
	EXT2_PIN10_SPI_SS_B = PB02
	EXT2_PIN10_GPIO     = PB02
	EXT2_PIN11_TWI_SDA  = PD08
	EXT2_PIN12_TWI_SCL  = PD09
	EXT2_PIN13_UART_RX  = PB17
	EXT2_PIN14_UART_TX  = PB16
	EXT2_PIN15_SPI_SS_A = PC06
	EXT2_PIN16_SPI_SDO  = PC04
	EXT2_PIN17_SPI_SDI  = PC07
	EXT2_PIN18_SPI_SCK  = PC05

	// Extension Header EXT3
	EXT3_PIN3_ADC_P     = PC02
	EXT3_PIN4_ADC_N     = PC03
	EXT3_PIN5_GPIO1     = PC01
	EXT3_PIN6_GPIO2     = PC10
	EXT3_PIN7_PWM_P     = PD10
	EXT3_PIN8_PWM_N     = PD11
	EXT3_PIN9_IRQ       = PC30
	EXT3_PIN9_GPIO      = PC30
	EXT3_PIN10_SPI_SS_B = PC31
	EXT3_PIN10_GPIO     = PC31
	EXT3_PIN11_TWI_SDA  = PD08
	EXT3_PIN12_TWI_SCL  = PD09
	EXT3_PIN13_UART_RX  = PC23
	EXT3_PIN14_UART_TX  = PC22
	EXT3_PIN15_SPI_SS_A = PC14
	EXT3_PIN16_SPI_SDO  = PC04
	EXT3_PIN17_SPI_SDI  = PC07
	EXT3_PIN18_SPI_SCK  = PC05

	// SD_CARD
	SD_CARD_MCDA0   = PB18
	SD_CARD_MCDA1   = PB19
	SD_CARD_MCDA2   = PB20
	SD_CARD_MCDA3   = PB21
	SD_CARD_MCCK    = PA21
	SD_CARD_MCCDA   = PA20
	SD_CARD_DETECT  = PD20
	SD_CARD_PROTECT = PD21

	// I2C
	I2C_SDA = PD08
	I2C_SCL = PD09

	// CAN
	CAN0_TX = PA22
	CAN0_RX = PA23

	CAN1_STANDBY = PC13
	CAN1_TX      = PB12
	CAN1_RX      = PB13

	CAN_STANDBY = CAN1_STANDBY
	CAN_TX      = CAN1_TX
	CAN_RX      = CAN1_RX

	// PDEC
	PDEC_PHASE_A = PC16
	PDEC_PHASE_B = PC17
	PDEC_INDEX   = PC18

	// PCC
	PCC_I2C_SDA    = PD08
	PCC_I2C_SCL    = PD09
	PCC_VSYNC_DEN1 = PA12
	PCC_HSYNC_DEN2 = PA13
	PCC_CLK        = PA14
	PCC_XCLK       = PA15
	PCC_DATA00     = PA16
	PCC_DATA01     = PA17
	PCC_DATA02     = PA18
	PCC_DATA03     = PA19
	PCC_DATA04     = PA20
	PCC_DATA05     = PA21
	PCC_DATA06     = PA22
	PCC_DATA07     = PA23
	PCC_DATA08     = PB14
	PCC_DATA09     = PB15
	PCC_RESET      = PC12
	PCC_PWDN       = PC11

	// Ethernet
	ETHERNET_TXCK  = PA14
	ETHERNET_TXEN  = PA17
	ETHERNET_TX0   = PA18
	ETHERNET_TX1   = PA19
	ETHERNET_RXER  = PA15
	ETHERNET_RX0   = PA13
	ETHERNET_RX1   = PA12
	ETHERNET_RXDV  = PC20
	ETHERNET_MDIO  = PC12
	ETHERNET_MDC   = PC11
	ETHERNET_INT   = PD12
	ETHERNET_RESET = PC21

	PIN_QT_BUTTON   = PA16
	PIN_BTN0        = PB31
	PIN_ETH_LED     = PC15
	PIN_LED0        = PC18
	PIN_ADC_DAC     = PA02
	PIN_VBUS_DETECT = PC00
	PIN_USB_ID      = PC19
)

// USBCDC pins
const (
	USBCDC_DM_PIN = PA24
	USBCDC_DP_PIN = PA25
)

// UART pins
const (
	// Extension Header EXT1
	UART_TX_PIN = PA04 // TX : SERCOM0/PAD[0]
	UART_RX_PIN = PA05 // RX : SERCOM0/PAD[1]

	// Extension Header EXT2
	UART2_TX_PIN = PB16 // TX : SERCOM5/PAD[0]
	UART2_RX_PIN = PB17 // RX : SERCOM5/PAD[1]

	// Extension Header EXT3
	UART3_TX_PIN = PC22 // TX : SERCOM1/PAD[0]
	UART3_RX_PIN = PC23 // RX : SERCOM1/PAD[1]

	// Virtual COM Port
	UART4_TX_PIN = PB25 // TX : SERCOM2/PAD[0]
	UART4_RX_PIN = PB24 // RX : SERCOM2/PAD[1]
)

// I2C pins
const (
	// Extension Header EXT1
	SDA0_PIN = PA22 // SDA: SERCOM3/PAD[0]
	SCL0_PIN = PA23 // SCL: SERCOM3/PAD[1]

	// Extension Header EXT2
	SDA1_PIN = PD08 // SDA: SERCOM7/PAD[0]
	SCL1_PIN = PD09 // SCL: SERCOM7/PAD[1]

	// Extension Header EXT3
	SDA2_PIN = PD08 // SDA: SERCOM7/PAD[0]
	SCL2_PIN = PD09 // SCL: SERCOM7/PAD[1]

	// Data Gateway Interface
	SDA_DGI_PIN = PD08 // SDA: SERCOM7/PAD[0]
	SCL_DGI_PIN = PD09 // SCL: SERCOM7/PAD[1]

	SDA_PIN = SDA0_PIN
	SCL_PIN = SCL0_PIN
)

// SPI pins
const (
	// Extension Header EXT1
	SPI0_SCK_PIN = PB26 // SCK: SERCOM4/PAD[1]
	SPI0_SDO_PIN = PB27 // SDO: SERCOM4/PAD[0]
	SPI0_SDI_PIN = PB29 // SDI: SERCOM4/PAD[3]
	SPI0_SS_PIN  = PB28 // SS : SERCOM4/PAD[2]

	// Extension Header EXT2
	SPI1_SCK_PIN = PC05 // SCK: SERCOM6/PAD[1]
	SPI1_SDO_PIN = PC04 // SDO: SERCOM6/PAD[0]
	SPI1_SDI_PIN = PC07 // SDI: SERCOM6/PAD[3]
	SPI1_SS_PIN  = PC06 // SS : SERCOM6/PAD[2]

	// Extension Header EXT3
	SPI2_SCK_PIN = PC05 // SCK: SERCOM6/PAD[1]
	SPI2_SDO_PIN = PC04 // SDO: SERCOM6/PAD[0]
	SPI2_SDI_PIN = PC07 // SDI: SERCOM6/PAD[3]
	SPI2_SS_PIN  = PC14 // SS : GPIO

	// Data Gateway Interface
	SPI_DGI_SCK_PIN = PC05 // SCK: SERCOM6/PAD[1]
	SPI_DGI_SDO_PIN = PC04 // SDO: SERCOM6/PAD[0]
	SPI_DGI_SDI_PIN = PC07 // SDI: SERCOM6/PAD[3]
	SPI_DGI_SS_PIN  = PD01 // SS : GPIO
)

// USB CDC identifiers
const (
	usb_STRING_PRODUCT      = "SAM E54 Xplained Pro"
	usb_STRING_MANUFACTURER = "Atmel"
)

var (
	usb_VID uint16 = 0x03EB
	usb_PID uint16 = 0x2404
)

// UART on the SAM E54 Xplained Pro
var (
	// Extension Header EXT1
	UART1  = &_UART1
	_UART1 = UART{
		Buffer: NewRingBuffer(),
		Bus:    sam.SERCOM0_USART_INT,
		SERCOM: 0,
	}

	// Extension Header EXT2
	UART2  = &_UART2
	_UART2 = UART{
		Buffer: NewRingBuffer(),
		Bus:    sam.SERCOM5_USART_INT,
		SERCOM: 5,
	}

	// Extension Header EXT3
	UART3  = &_UART3
	_UART3 = UART{
		Buffer: NewRingBuffer(),
		Bus:    sam.SERCOM1_USART_INT,
		SERCOM: 1,
	}

	// EDBG Virtual COM Port
	UART4  = &_UART4
	_UART4 = UART{
		Buffer: NewRingBuffer(),
		Bus:    sam.SERCOM2_USART_INT,
		SERCOM: 2,
	}
)

func init() {
	UART1.Interrupt = interrupt.New(sam.IRQ_SERCOM0_2, _UART1.handleInterrupt)
	UART2.Interrupt = interrupt.New(sam.IRQ_SERCOM5_2, _UART2.handleInterrupt)
	UART3.Interrupt = interrupt.New(sam.IRQ_SERCOM1_2, _UART3.handleInterrupt)
	UART4.Interrupt = interrupt.New(sam.IRQ_SERCOM2_2, _UART4.handleInterrupt)
}

// I2C on the SAM E54 Xplained Pro
var (
	// Extension Header EXT1
	I2C0 = I2C{
		Bus:    sam.SERCOM3_I2CM,
		SERCOM: 3,
	}

	// Extension Header EXT2
	I2C1 = I2C{
		Bus:    sam.SERCOM7_I2CM,
		SERCOM: 7,
	}

	// Extension Header EXT3
	I2C2 = I2C{
		Bus:    sam.SERCOM7_I2CM,
		SERCOM: 7,
	}

	// Data Gateway Interface
	I2C3 = I2C{
		Bus:    sam.SERCOM7_I2CM,
		SERCOM: 7,
	}
)

// SPI on the SAM E54 Xplained Pro
var (
	// Extension Header EXT1
	SPI0 = SPI{
		Bus:    sam.SERCOM4_SPIM,
		SERCOM: 4,
	}

	// Extension Header EXT2
	SPI1 = SPI{
		Bus:    sam.SERCOM6_SPIM,
		SERCOM: 6,
	}

	// Extension Header EXT3
	SPI2 = SPI{
		Bus:    sam.SERCOM6_SPIM,
		SERCOM: 6,
	}

	// Data Gateway Interface
	SPI3 = SPI{
		Bus:    sam.SERCOM6_SPIM,
		SERCOM: 6,
	}
)

// CAN on the SAM E54 Xplained Pro
var (
	CAN0 = CAN{
		Bus: sam.CAN0,
	}

	CAN1 = CAN{
		Bus: sam.CAN1,
	}
)
