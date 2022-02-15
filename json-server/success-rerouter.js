const getSuccessRoute = (req) => {
    return req.headers?.cookie?.match(/success-route=(?<failRoute>\w+);/)
        ?.groups.failRoute;
};

module.exports = (req, res, next) => {
    if (["POST", "PUT", "PATCH"].includes(req.method)) {
        const successRoute = getSuccessRoute(req);

        if (successRoute) {
            req.method = "GET";
            req.url = `/successes/${successRoute}`;
            res.status(200);
        }
    }
    next();
};
