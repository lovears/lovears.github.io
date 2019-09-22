---
layout:     post 
title:      "【线程池】java thread"
subtitle:   "java中的线程和线程池"
date:       2019-09-22
author:     "老回路"
header-img: "img/home-bg.png"
tags:
    - thread
    - threadpool
    - 线程池
---

# Java中的线程和多线程

## 基本概念
线程是操作系统能够进行运算调度的最小单位。它被包含在进程中，是进程中的实际运作单位。一条线程指的是进程中一个单一顺序控制流，一个进程中可以并发多个线程，每个线程并行执行不同的任务。  
线程是独立调度和分派的基本单位。  
统一进程中的多个线程将共享该进程中的所有资源，如虚拟地址空间，文件描述符和信号处理等。单统一进程中的多条线程有各自的调用栈，自己的寄存器环境，自己的线程本地存储。

## 线程的状态
在Java中的线程，具有物种基本状态:
- 新建 (New)
- 就绪 (Runnable)
- 运行 (Running))
- 阻塞 (Blocked)
- 死亡 (Dead)

### **状态流转图：**
![statemachine](/img/in-post/java-thread-state-machine.png)

## 线程通信
- wait()/notify() 等待通知
- ReentrantLock类加锁的线程的Condition类的await()/signal()/signalAll()
- 通过管道进行线程间通信：1）字节流；2）字符流

## Java中的线程池
- FixedThreadPool: 
只有核心线程，且数量固定，没有非核心线程。keepAliveTime设置为0L，代表多余的线程会被立即终止。任务队列采用了无界的阻塞队列LinkedBlockingQueue（容量默认为Integer.MAX_VALUE）。

```java
/**
     * Creates a thread pool that reuses a fixed number of threads
     * operating off a shared unbounded queue, using the provided
     * ThreadFactory to create new threads when needed.  At any point,
     * at most {@code nThreads} threads will be active processing
     * tasks.  If additional tasks are submitted when all threads are
     * active, they will wait in the queue until a thread is
     * available.  If any thread terminates due to a failure during
     * execution prior to shutdown, a new one will take its place if
     * needed to execute subsequent tasks.  The threads in the pool will
     * exist until it is explicitly {@link ExecutorService#shutdown
     * shutdown}.
     *
     * @param nThreads the number of threads in the pool
     * @param threadFactory the factory to use when creating new threads
     * @return the newly created thread pool
     * @throws NullPointerException if threadFactory is null
     * @throws IllegalArgumentException if {@code nThreads <= 0}
     */
    public static ExecutorService newFixedThreadPool(int nThreads, ThreadFactory threadFactory) {
        return new ThreadPoolExecutor(nThreads, nThreads,
                                      0L, TimeUnit.MILLISECONDS,
                                      new LinkedBlockingQueue<Runnable>(),
                                      threadFactory);
    }
```

- CachedThreadPool:
创建之初里面一个线程都没有，当execute方法或submit方法向线程池提交任务时，会自动新建线程；如果线程池中有空余线程，则不会新建；这种线程池一般最多情况可以容纳几万个线程，里面的线程空余60s会被回收

```java
/**
     * Creates a thread pool that creates new threads as needed, but
     * will reuse previously constructed threads when they are
     * available.  These pools will typically improve the performance
     * of programs that execute many short-lived asynchronous tasks.
     * Calls to {@code execute} will reuse previously constructed
     * threads if available. If no existing thread is available, a new
     * thread will be created and added to the pool. Threads that have
     * not been used for sixty seconds are terminated and removed from
     * the cache. Thus, a pool that remains idle for long enough will
     * not consume any resources. Note that pools with similar
     * properties but different details (for example, timeout parameters)
     * may be created using {@link ThreadPoolExecutor} constructors.
     *
     * @return the newly created thread pool
     */
    public static ExecutorService newCachedThreadPool() {
        return new ThreadPoolExecutor(0, Integer.MAX_VALUE,
                                      60L, TimeUnit.SECONDS,
                                      new SynchronousQueue<Runnable>());
    }
```

- SingleThreadPool:
池中只有一个线程，多余任务将排队；作用是保证任务的顺序执行；如果单独的线程因为失败而终止，将会使用新的线程代替

```java
 /**
     * Creates a single-threaded executor that can schedule commands
     * to run after a given delay, or to execute periodically.
     * (Note however that if this single
     * thread terminates due to a failure during execution prior to
     * shutdown, a new one will take its place if needed to execute
     * subsequent tasks.)  Tasks are guaranteed to execute
     * sequentially, and no more than one task will be active at any
     * given time. Unlike the otherwise equivalent
     * {@code newScheduledThreadPool(1)} the returned executor is
     * guaranteed not to be reconfigurable to use additional threads.
     * @return the newly created scheduled executor
     */
    public static ScheduledExecutorService newSingleThreadScheduledExecutor() {
        return new DelegatedScheduledExecutorService
            (new ScheduledThreadPoolExecutor(1));
    }
```

- ScheduledThreadpool
可以给定延迟时间或者定期执行

```java
   /**
     * Creates a thread pool that can schedule commands to run after a
     * given delay, or to execute periodically.
     * @param corePoolSize the number of threads to keep in the pool,
     * even if they are idle
     * @param threadFactory the factory to use when the executor
     * creates a new thread
     * @return a newly created scheduled thread pool
     * @throws IllegalArgumentException if {@code corePoolSize < 0}
     * @throws NullPointerException if threadFactory is null
     */
    public static ScheduledExecutorService newScheduledThreadPool(
            int corePoolSize, ThreadFactory threadFactory) {
        return new ScheduledThreadPoolExecutor(corePoolSize, threadFactory);
    }
```

## Java线程池拒绝策略
- AbortPolicy： 直接抛出异常
- CallerRunsPolicy： 调用者线程执行
- DiscardOldestPolicy：丢弃队列中最老的任务
- DiscardPolicy： 直接丢弃

## 线程池的实现原理
- 一个工作队列`BlockingQueue<Runnable> workQueue;` 用来保存所有的任务
- 任务线程集合` HashSet<Worker> workers`
- 其他参数： 线程池大小，线程有效时间等。

### 任务提交到执行

提交任务，到工作队列中保存，`Worker`中有一个`while(true)` 一直从 workQueue中取任务，如果存在就执行，不存在就一直循环