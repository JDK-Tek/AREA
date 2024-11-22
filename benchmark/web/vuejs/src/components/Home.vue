<template>
    <div class='Menu'>
        <Menu/>
    </div>
    <div class="div-techno">
        <img :src="reactJSLogo" class="img-techno" alt="React"/>
        <img :src="HTMLLogo" class="img-techno" alt="React"/>
        <img :src="VueJSLogo" class="img-techno" alt="React"/>
    </div>
  </template>

<style>
.Menu {
    text-align: center;
    margin-top: 1.3rem;
}

.div-techno {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-top: 20px;
}

.img-techno {
  width: 20vh;
  height: 20vh;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  object-fit: cover;
  display: block;
  margin: auto;
  padding: 20px;
}
</style>

<script>
    import reactJS from '@/assets/web/reactjs.png';
    import HTML from '@/assets/web/html.png';
    import VueJS from '@/assets/web/vuejs.png';

    import Django from '@/assets/backend/django.png';
    import GO from '@/assets/backend/go.png';
    import NestJS from '@/assets/backend/nestjs.jpg';

    import MariaDB from '@/assets/database/mariadb.jpeg';
    import MySQL from '@/assets/database/mysql.png';
    import PostGres from '@/assets/database/postgres.png';
    import SQLite from '@/assets/database/sqlite.png';

    import Flutter from '@/assets/mobile/flutter.png';
    import Kotlin from '@/assets/mobile/kotline.png';
    import ReactNative from '@/assets/mobile/reactnative.png';
    
    import Menu from './Menu.vue'

    export default {
        name: 'Home',
        components: {
            Menu
        },
        data() {
            return {
                reactJSLogo: reactJS,
                HTMLLogo: HTML,
                VueJSLogo: VueJS,
                DjangoLogo: Django,
                GOLogo: GO,
                NestJSLogo: NestJS,
                MariaDBLogo: MariaDB,
                MySQLLogo: MySQL,
                PostGresLogo: PostGres,
                SQLiteLogo: SQLite,
                FlutterLogo: Flutter,
                KotlinLogo: Kotlin,
                ReactNativeLogo: ReactNative
            };
        }
    };
  </script>
