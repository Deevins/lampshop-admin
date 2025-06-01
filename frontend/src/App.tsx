import React from "react";
import {Routes, Route, Navigate} from "react-router-dom";
import {AuthProvider} from "./contexts/AuthContext";
import ProtectedRoute from "./routes/ProtectedRoute";
import Layout from "./components/Layout/Layout";
import ProductList from "./Pages/Products/ProductList.tsx";
import ProductForm from "./Pages/Products/ProductForm.tsx";
import Login from "./Pages/Login/Login.tsx";
import OrderList from "./Pages/Orders/OrderList.tsx";


const App: React.FC = () => {
    return (
        <AuthProvider>
            <Routes>
                {/* Public Route */}
                <Route path="/login" element={<Login/>}/>

                {/* All other routes are protected */}
                <Route
                    path="/"
                    element={
                        <ProtectedRoute>
                            <Layout/>
                        </ProtectedRoute>
                    }
                >
                    {/* Inside Layout, nested routes: */}
                    <Route index element={<Navigate to="/products"/>}/>
                    <Route path="products" element={<ProductList/>}/>
                    <Route path="products/add" element={<ProductForm/>}/>
                    <Route path="products/edit/:id" element={<ProductForm/>}/>
                    <Route path="orders" element={<OrderList/>}/>
                </Route>

                {/* Catch-all redirect to /login */}
                {/*<Route path="*" element={<Navigate to="/login"/>}/>*/}
            </Routes>
        </AuthProvider>
    );
};

export default App;
