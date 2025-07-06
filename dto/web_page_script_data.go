package dto

type SDProduct struct {
	Id string `json:"id"`
}

type PLPLink struct {
	Url string `json:"url"`
}

type PageSeo struct {
	PlpLinks []PLPLink `json:"plpLinks,omitempty"`
}

type NavItem struct {
	Href string `json:"href"`
}

type NavigationData struct {
	NavItems []*NavItem `json:"navItems"`
}

type PageProps struct {
	PageType       string          `json:"pageType"`
	Products       []SDProduct     `json:"products"`
	Layouts        []*Layout       `json:"layouts,omitempty"`
	PageSeo        *PageSeo        `json:"pageSeo,omitempty"`
	PathName       string          `json:"pathname"`
	NavigationData *NavigationData `json:"navigationData,omitempty"`
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

type CTA struct {
	RelativeUrl string `json:"url"`
}

type Content struct {
	CTAs []*CTA `json:"ctas"`
}

type Layout struct {
	Contents []*Content `json:"contents"`
}

type WebPageScriptData struct {
	Props *Props `json:"props"`
}
