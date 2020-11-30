import http from "k6/http";
import { check } from "k6";
export default function () {
    const res = http.get("http://localhost:4000/protected", {
        redirects: 0,
    });
    check(res, {
        "response code was 303": (res) => res.status == 303,
    });
}
