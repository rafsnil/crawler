package constant

const (
	PRODUCT_URL_IDENTIFIER       = "product-card-image-link" // a -> attr: data-testid
	BREADCRUMB_IDENTIFIER        = "breadcrumbs-desktop"     // ol -> attr: data-auto-id (Only if the visited url is a product page)
	LOOKBOOK_SSR_DATA_IDENTIFIER = "lookbook-microfrontend"  // script -> attr: data-mf-id
	ALL_COORDINATES_IDENTIFIER   = "styles-carousel"         // div -> attr: data-testid

	SIZE_CHART_IDENTIFIER     = "garment-measurement"
	NO_SIZE_CHART             = "size-chart-NA"
	SCRIPT_DATA_IDENTIFIER    = "__NEXT_DATA__"
	PAGE_TYPE_PRODUCT_LISTING = "ProductListingPage"
	PAGE_TYPE_LANDING         = "GenderLandingPage"
	VISITED                   = "visited.txt" // File to store visited URLs
	QUEUE                     = "queue.txt"   // File to store URLs to be crawled
	DATA_LIMIT                = 250
)
