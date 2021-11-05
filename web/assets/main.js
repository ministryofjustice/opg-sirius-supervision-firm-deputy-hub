import './main.scss';
import GOVUKFrontend from 'govuk-frontend/govuk/all.js';
import MojBannerAutoHide from './javascript/moj-banner-auto-hide';

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector('.app-main-class'));
