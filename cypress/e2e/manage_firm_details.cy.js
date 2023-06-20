describe("Manage Firm Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        it("has a link from the main firms page", () => {
            cy.visit("/supervision/deputies/firm/1");
            cy.get(":nth-child(2) > .govuk-button")
                .should("contain", "Manage firm details")
                .click();
            cy.url().should("include", "firm/1/manage-firm-details");
            cy.get("#manage-firm-title").should("be.visible");
        });

        // it("has a cancel link which returns to the main firms page", () => {
        //     cy.visit("/supervision/deputies/firm/1/manage-firm-details");
        //     cy.get(".govuk-link").should("contain", "Cancel").click();
        //     cy.url().should("not.include", "firm/1/manage-firm-details");
        //     cy.get(".govuk-heading-l").should("contain", "Firm details");
        // });
    });

    describe("Edit Firm Details Form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/firm/1");
            cy.get(":nth-child(2) > .govuk-button")
                .should("contain", "Manage firm details")
                .click();
        });

        it("should autopopulate the existing firm data", () => {
            cy.get("#f-firm-name").should("have.value", "Trustworthy Firm Inc");
            cy.get("#address-line-1").should("have.value", "221 Baker Street");
            cy.get("#address-line-2").should("have.value", "Covent Garden");
            cy.get("#town").should("have.value", "London");
            cy.get("#county").should("have.value", "Buckinghamshire");
            cy.get("#postcode").should("have.value", "B10 1FG");
        });

        it("should allow for inputs to be amended or added", () => {
            cy.get("#address-line-1").clear();
            cy.get("#address-line-1").type("Amended address line 1");
            cy.get("#address-line-2").type(" Road");
            cy.get("#address-line-3").type("West London");
        });

        it("allows me to successfully submit the form", () => {
            cy.setCookie("success-route", "manageFirm");
            cy.get("#address-line-1").clear();
            cy.get("#address-line-1").type("Amended address line 1");
            cy.get("#manage-firm-details-submit-btn").click();
            cy.url().should("not.include", "firm/1/manage-firm-details");
        });

        it("will show a validation error if form submitted with empty fields", () => {
            cy.setCookie("fail-route", "manage-firm-details-empty");
            cy.get("[data-cy=submit-manage-firm-details-form-btn]").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-list > li").should(
                "contain",
                "The firm name is required and can't be empty"
            );
        });
        it("will show a validation error if form submitted with empty fields", () => {
            cy.setCookie("fail-route", "manage-firm-details-too-long");
            cy.get("[data-cy=submit-manage-firm-details-form-btn]").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-list > li").should(
                "contain",
                "The firm name must be 255 characters or fewer"
            );
        });
    });
});
