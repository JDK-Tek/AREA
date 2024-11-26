<template>
    <div class='Menu'>
        <Menu @category-selected="changeCategory"/>
    </div>
    <div class="div-techno">
        <img v-for="image in filteredImages" :key="image" :src="image" class="img-techno" :alt="selectedCategory"/>
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
                selectedCategory: 'frontend',
                images: {
                    Web: [reactJS, HTML, VueJS],
                    Backend: [Django, GO, NestJS],
                    Database: [MariaDB, MySQL, PostGres, SQLite],
                    Mobile: [Flutter, Kotlin, ReactNative]
                }
            };
        },
        computed: {
            filteredImages() {
            return this.images[this.selectedCategory];
            }
        },
        methods: {
            changeCategory(category) {
            this.selectedCategory = category;
            }
        }
    };
  </script>
