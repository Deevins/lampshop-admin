import React from "react";
import {Routes, Route, Navigate} from "react-router-dom";
import ProtectedRoute from "./routes/ProtectedRoute";
import ProductList from "./Pages/Products/ProductList.tsx";
import ProductForm from "./Pages/Products/ProductForm.tsx";
import Login from "./Pages/Login/Login.tsx";
import OrderList from "./Pages/Orders/OrderList.tsx";
import Layout from "./components/Layout/Layout.tsx";


const App: React.FC = () => {
    return (
        <Routes>
            <Route path="/login" element={<Login/>}/>

            <Route
                path="/"
                element={
                    <ProtectedRoute>
                        <Layout/>
                    </ProtectedRoute>
                }
            >
                <Route path="/" element={<Navigate to="/products" replace/>}/>
                <Route
                    path="products"
                    element={
                        <ProtectedRoute>
                            <ProductList/>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="products/add"
                    element={
                        <ProtectedRoute>
                            <ProductForm/>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="products/edit/:id"
                    element={
                        <ProtectedRoute>
                            <ProductForm/>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="orders"
                    element={
                        <ProtectedRoute>
                            <OrderList/>
                        </ProtectedRoute>
                    }
                />
            </Route>


            <Route path="*" element={<Navigate to="/login" replace/>}/>
        </Routes>
    );
};

export default App;