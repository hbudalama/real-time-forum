// router.js
document.addEventListener("DOMContentLoaded", () => {
    const navigateTo = url => {
        history.pushState(null, null, url);
        router();
    };

    const router = async () => {
        const routes = [
            { path: "/", view: () => fetch("/api/home").then(res => res.text()) },
            { path: "/login", view: () => fetch("/api/login").then(res => res.text()) },
            { path: "/posts", view: () => fetch("/api/posts").then(res => res.text()) },
            // { path: "/chat", view: () => fetch("/api/chat").then(res => res.text()) },
        ];

        const potentialMatches = routes.map(route => {
            return {
                route,
                isMatch: location.pathname === route.path
            };
        });

        let match = potentialMatches.find(potentialMatch => potentialMatch.isMatch);

        if (!match) {
            match = {
                route: routes[0],
                isMatch: true
            };
        }

        const view = await match.route.view();

        document.querySelector("#main-content").innerHTML = view;
    };

    window.addEventListener("popstate", router);

    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    router();
});