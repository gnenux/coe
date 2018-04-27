(function () {
    function addLoadListener(fn) {
        if (typeof window.addEventListener != "undefined") {
            window.addEventListener("load", fn, false)
        } else {
            if (typeof document.addEventListener) {
                document.addEventListener("load", fn, false)
            } else {
                if (typeof window.attachEvent != "undefined") {
                    window.attachEvent("load", fn)
                } else {
                    var oldfn = window.onload;
                    if (typeof window.onload != "function") {
                        window.onload = fn
                    } else {
                        window.onload = function () {
                            oldfn();
                            fn()
                        }
                    }
                }
            }
        }
    }


    addLoadListener(function () {

        function sendData(form) {
            var url = form.action + "?";
            var fd = new FormData(form);
            var searchStr = fd.get("search");
            searchKeys = searchStr.split(" ");
            for (let index = 0; index < searchKeys.length; index++) {
                const element = searchKeys[index];
                url = url + "key=" + searchKeys[index];
                if (index < searchKeys.length - 1) {
                    url = url + "&";
                }
            }
            window.location.href = url;
        }

        var navForm = document.getElementById("nav-search-form");
        if (navForm != null && typeof (navForm) != "undefined") {
            navForm.addEventListener("submit", function (event) {
                event.preventDefault();
                sendData(navForm);
            });

        }

        var mainForm = document.getElementById("main-search-form");
        if (mainForm != null && typeof (mainForm) != "undefined") {
            mainForm.addEventListener("submit", function (event) {
                event.preventDefault();
                sendData(mainForm);
            });
        }


    });

})()