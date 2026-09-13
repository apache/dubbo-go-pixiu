import React from 'react';
import {createBrowserRouter,Navigate,Outlet} from 'react-router-dom';
import {App} from './App';
import {Login} from '../features/profile/Login';
import {readSession} from '../services/session';
function RequireAuth(){return readSession()?<Outlet/>:<Navigate to="/login" replace/>}
export const router=createBrowserRouter([{path:'/login',element:<Login/>},{element:<RequireAuth/>,children:[{path:'/',element:<Navigate to="/gateway/overview" replace/>},{path:'/gateway/overview',element:<App/>},{path:'/gateway/*',element:<App/>}]},{path:'*',element:<Navigate to="/gateway/overview" replace/>}]);
