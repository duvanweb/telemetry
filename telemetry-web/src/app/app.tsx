import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Dashboard } from "@/presentation/pages/dashboard";

// App — router with a single route rendering the dashboard.
export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
      </Routes>
    </BrowserRouter>
  );
}
