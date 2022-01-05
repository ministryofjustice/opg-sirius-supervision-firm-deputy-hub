describe("Manage PII Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Form", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/firm/1/manage-pii-details");
        });

        it("will show a validation error if form submitted with empty fields", () => {
            cy.setCookie("fail-route", "pii-details");
            cy.get("[data-cy=submit-manage-pii-details-form-btn]").click();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );

            cy.get(".govuk-error-summary").should(
                "contain",
                "The PII expiry date is required and can't be empty"
            ); 
            
            cy.get(".govuk-error-summary").should(
                "contain",
                "The PII received date is required and can't be empty"
            );

            cy.get(".govuk-error-summary").should(
                "contain",
                "The PII amount is required and can't be empty"
            );
        });

        it("will autofill the form fields", () => {
            cy.get("#f-pii-received").should(
                "have.value",
                "2000-12-20"
            );

            cy.get("#f-pii-expiry").should(
                "have.value",
                "2020-12-01"
            );

            cy.get("#f-pii-amount").should(
                "have.value",
                "1000"
            );
        });
    });
})