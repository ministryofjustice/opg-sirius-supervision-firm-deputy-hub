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

        const expected = ["Workflow", "Supervision", "LPA", "Admin", "Logout"];

        it("has working nav links within header banner", () => {
            cy.get(".moj-header__navigation-list")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                    const $linkName = expected[index].toLowerCase();
                    cy.wrap($el)
                        .find("a")
                        .should("have.attr", "href")
                        .and("contain", `/${$linkName}`);
                });
        });
    });

    describe("Firm Details Header", () => {
        it("should show the firm name", () => {
            cy.get(".govuk-grid-column-full > .govuk-heading-m").should(
                "contain",
                "Trustworthy Firm Inc"
            );
        });

        it("should show the firm number", () => {
            cy.get(".govuk-caption-m.govuk-\\!-margin-bottom-0").should(
                "contain",
                "100004"
            );
        });

        it("should have space for the executive Case manager to be added in the future", () => {
            cy.get(".govuk-\\!-margin-bottom-2").should(
                "contain",
                "Executive Case Manager"
            );
        });
    });

    describe("Firm Details Navigation", () => {
        it("has a link for the current page", () => {
            cy.get(":nth-child(1) > .moj-sub-navigation__link").should(
                "contain",
                "Firm details"
            );
        });
    });

    describe("Firm Details Body", () => {
        it("should show the firm name", () => {
            cy.get(
                "#team-details > :nth-child(1) > .govuk-summary-list__key"
            ).should("contain", "Firm name");
            cy.get(
                "#team-details > :nth-child(1) > .govuk-summary-list__value"
            ).should("contain", "Trustworthy Firm Inc");
        });

        it("should show the firm address", () => {
            cy.get(
                "#team-details > :nth-child(2) > .govuk-summary-list__key"
            ).should("contain", "Main address");
            cy.get(
                "#team-details > :nth-child(2) > .govuk-summary-list__value"
            ).should("contain", "221 Baker Street");
        });

        it("should show the phone number", () => {
            cy.get(
                "#team-details > :nth-child(3) > .govuk-summary-list__key"
            ).should("contain", "Telephone");
            cy.get(
                "#team-details > :nth-child(3) > .govuk-summary-list__value"
            ).should("contain", "333222111");
        });

        it("should show the email address with a mail to link", () => {
            cy.get(
                "#team-details > :nth-child(4) > .govuk-summary-list__key"
            ).should("contain", "Email");
            cy.get(
                "#team-details > :nth-child(4) > .govuk-summary-list__value"
            ).should("contain", "trusty@firm.com");
        });
    });

    describe("Firm PII Details Body", () => {
        it("should show the PII expiry date", () => {
            cy.get(
                ":nth-child(3) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__value"
            ).should("contain", "20/12/2000");
            cy.get(
                ":nth-child(3) > .govuk-summary-list > :nth-child(2) > .govuk-summary-list__value"
            ).should("contain", "01/12/2020");
            cy.get(
                ":nth-child(3) > .govuk-summary-list > :nth-child(3) > .govuk-summary-list__value"
            ).should("contain", "1,000");
            cy.get(
                ":nth-child(3) > .govuk-summary-list > :nth-child(4) > .govuk-summary-list__value"
            ).should("contain", "15/11/2000");
        });
    });

    describe("Navigation", () => {
        it("should navigate to the 'Manage PIIS contact details' page", () => {
            cy.get("[data-cy=manage-pii-details-btn]").click();
            cy.contains(
                ".govuk-heading-l",
                "Manage professional indemnity insurance"
            );
        });
    });

    describe("Footer", () => {
        it("the footer should contain a link to the open government licence", () => {
            cy.get(
                ".govuk-footer__licence-description > .govuk-footer__link"
            ).should(
                "have.attr",
                "href",
                "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/"
            );
        });

        it("the nav link should contain the crown copyright logo", () => {
            cy.get(".govuk-footer__copyright-logo").should(
                "have.attr",
                "href",
                "https://www.nationalarchives.gov.uk/information-management/re-using-public-sector-information/uk-government-licensing-framework/crown-copyright/"
            );
        });
    });
});
