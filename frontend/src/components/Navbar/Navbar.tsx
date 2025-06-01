import React from "react";
import { NavLink } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import styles from "./Navbar.module.scss";

const Navbar: React.FC = () => {
    const { logout } = useAuth();

    return (
        <div className={styles.navbar}>
            <div className={styles["nav-links"]}>
                <NavLink
                    to="/products"
                    className={({ isActive }) =>
                        isActive ? `${styles["nav-link"]} ${styles["active"]}` : styles["nav-link"]
                    }
                >
                    Товары
                </NavLink>
                <NavLink
                    to="/orders"
                    className={({ isActive }) =>
                        isActive ? `${styles["nav-link"]} ${styles["active"]}` : styles["nav-link"]
                    }
                >
                    Заказы
                </NavLink>
            </div>
            <button className={styles["logout-button"]} onClick={logout}>
                Выход
            </button>
        </div>
    );
};

export default Navbar;
