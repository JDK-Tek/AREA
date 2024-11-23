import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { HelloController } from './hello/hello.controller';
import { ApiController } from './api/api.controller';

@Module({
  imports: [],
  controllers: [AppController, HelloController, ApiController],
  providers: [AppService],
})
export class AppModule {}
