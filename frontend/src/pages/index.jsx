// pages/index.jsx
export default function Homepage() {
  return (
    <div className="flex min-h-screen flex-col">
      {/* 1a) Nav lives in _app.js, so main just handles the hero */}
      <main className="from-primary-500 to-accent flex flex-1 flex-col items-center justify-center bg-gradient-to-br px-4 text-center text-black">
        <h1 className="text-6xl font-extrabold tracking-tight sm:text-7xl">
          Welcome to Nexus Local
        </h1>
        <p className="mt-4 max-w-2xl text-lg opacity-90">
          Your one-stop shop for unique local finds.
        </p>
        <button
          onClick={() => (window.location.href = "/listing")}
          className="text-primary-600 hover:shadow-3xl mt-8 rounded-full bg-white px-8 py-4 text-lg font-medium shadow-2xl transition"
        >
          Browse Products
        </button>
      </main>
    </div>
  );
}
