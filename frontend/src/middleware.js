// // pages/admin/add-item.jsx
// import { parse } from "cookie";
//
// export async function getServerSideProps({ req }) {
//   const cookies = parse(req.headers.cookie || "");
//   if (!cookies.id_token) {
//     return {
//       redirect: {
//         destination: "http://localhost:8080/login?redirect_uri=/admin/add-item",
//         permanent: false,
//       },
//     };
//   }
//
//   // Optionally: verify the token by calling your Go backend
//   // const resp = await fetch('http://localhost:8080/verify', {
//   //   headers: { Cookie: req.headers.cookie }
//   // });
//   // if (resp.status !== 200) { ...redirect again }
//
//   return { props: {} };
// }
//
// export default function AddItemPage() {
//   // your existing form component…
// }
//
// middleware.js
import { NextResponse } from "next/server";

export function middleware(req) {
  const { pathname } = req.nextUrl;

  if (pathname.startsWith("/admin")) {
    const idToken = req.cookies.get("id_token");
    if (!idToken) {
      const loginUrl = new URL("http://localhost:8080/login");
      loginUrl.searchParams.set("redirect_uri", req.nextUrl.href);
      return NextResponse.redirect(loginUrl);
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*"],
};
