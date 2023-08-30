import "cypress-axe";
import { TerminalLog } from "../support/e2e";
import navTabs from "../fixtures/navigation.json"

// describe("Accessibility", () => {
//     before(() => {
//         cy.visit("/supervision/deputies/firm/1");
//         // cy.url().should('contain', '/deputies/firm/1')
//         cy.injectAxe();
//     });
//
//     it("Should have no accessibility violations",() => {
//         cy.checkA11y(null, null, TerminalLog)
//     });
// });

describe("Accessibility", () => {
    navTabs.forEach(([page, url]) =>
        it(`should render ${page} page accessibly`, () => {
            cy.visit(url);
            cy.injectAxe();
            cy.checkA11y(null, null, TerminalLog);
        })
    )
});
