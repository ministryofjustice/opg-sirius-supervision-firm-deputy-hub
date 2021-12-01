describe("Firm Deputy Hub", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/firm/1");
    });

    describe("Header", () => {
        it("shows opg sirius within banner", () => {
            cy.contains(".moj-header__link", "OPG");
            cy.contains(".moj-header__link", "Sirius");
        });

        const expected = [
            "Supervision",
            "LPA",
            "Admin",
            "Logout",
        ];

        it("has working nav links within header banner", () => {
            cy.get(".moj-header__navigation-list")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                    let $linkName = (expected[index].toLowerCase());
                    cy.wrap($el).find('a').should("have.attr", "href").and("contain", `/${$linkName}`);

                });
        });
    });

    describe("Firm Details", () => {
        it("should show the firm name", () => {
            cy.get(".govuk-grid-column-full > .govuk-heading-m").should("contain",  "Trustworthy Firm Inc")
        })

        it("should show the firm number", () => {
            cy.get(".govuk-caption-m.govuk-\\!-margin-bottom-0").should("contain",  "100004")
        })

        it("should have space for the executive Case manager to be added in the future", () => {
            cy.get(".govuk-\\!-margin-bottom-2").should("contain",  "Executive Case Manager")
        })
    });

    describe("Footer", () => {
        it("the footer should contain a link to the open government licence", () => {
            cy.get(".govuk-footer__licence-description > .govuk-footer__link").should("have.attr", "href", "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/")
        })

        it("the nav link should contain the crown copyright logo", () => {
            cy.get(".govuk-footer__copyright-logo").should("have.attr", "href", "https://www.nationalarchives.gov.uk/information-management/re-using-public-sector-information/uk-government-licensing-framework/crown-copyright/")
        })
    });
});