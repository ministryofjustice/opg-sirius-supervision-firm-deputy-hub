import "govuk-frontend/dist/govuk/all.mjs";
import MojBannerAutoHide from "./javascript/moj-banner-auto-hide";
import accessibleAutocomplete from "accessible-autocomplete";
import "opg-sirius-header/sirius-header.js";

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector(".app-main-class"));

if (document.querySelector("#select-ecm")) {
    accessibleAutocomplete.enhanceSelectElement({
        selectElement: document.querySelector("#select-ecm"),
        defaultValue: "",
    });
}

if (document.querySelector("#f-back-button")) {
    document.getElementById("f-back-button").onclick = function (e) {
        e.preventDefault();
        history.go(parseInt(sessionStorage.getItem("backIndex")));
    }
}

function onHomePage() {
    const homePageUrlRegex = new RegExp('^\\/(supervision/deputies/firm\\/)?\\d+\\/*$');
    return homePageUrlRegex.test(location.pathname);
}

function storeBackSessionVars(backIndex, href) {
    if (backIndex !== null && location.href === href) {
        sessionStorage.setItem("backIndex", (parseInt(backIndex) - 1).toString());
    }

    if (backIndex === null || href === null || location.href !== href || onHomePage()) {
        sessionStorage.setItem("backIndex", "-1");
        sessionStorage.setItem("href", location.href);
    }
}

storeBackSessionVars(sessionStorage.getItem("backIndex"), sessionStorage.getItem("href"));
