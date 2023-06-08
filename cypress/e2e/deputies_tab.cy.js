describe("Deputies Tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("List of deputies", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/firm/1/deputies");
        });

        it("shows the page title", () => {
            cy.get(".govuk-heading-l").contains("Deputies");
        });

        it("have depuites name as hyper links", () => {
            cy.get(":nth-child(1) > .client_name_ref > .govuk-link").contains(
                "pro Deputy"
            );
            cy.get(":nth-child(1) > .client_name_ref > .govuk-link").should(
                "have.attr",
                "href"
            );
        });

        it("shows a dash if no assurance visit", () => {
            cy.get(
                ".govuk-table__body > :nth-child(2) > :nth-child(5)"
            ).contains("-");
        });

        it("shows assurance visit data", () => {
            cy.get(":nth-child(1) > .visit_type").contains(
                "26/05/2023"
            );
            cy.get(":nth-child(1) > .visit_type > .secondary").contains(
                "Green"
            );
        });

        it("shows a dash if no ECM", () => {
            cy.get(
                ".govuk-table__body > :nth-child(2) > :nth-child(6)"
            ).contains("-");
        });

        it("has a Firm details tab that is clickable", () => {
            cy.get(":nth-child(1) > .moj-sub-navigation__link").click();
            cy.url().should(
                "not.eq",
                "http://localhost:8887/supervision/deputies/firm/1/deputies"
            );
        });

        it("displays the deputies details", () => {
            const expected = [
                "pro Deputy",
                "The Town",
                "22",
                "3",
                "PROTeam2 User1",
            ];
            cy.get(".govuk-table__body > :first-child").each(($elm, index) => {
                expect($elm.text()).to.contain(expected[index]);
            });
        });
    });
});
