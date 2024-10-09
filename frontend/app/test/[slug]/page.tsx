"use client"
import { useEffect, useRef } from 'react';
import Hls from 'hls.js';
import { useParams } from 'next/navigation';
export default function Page() {
  const { slug } = useParams();
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const hlsUrl = `https://publicstream.s3.ap-south-1.amazonaws.com/${slug}/playlist.m3u8`;
  console.log(hlsUrl);

  useEffect(() => {
    const video = videoRef.current;

    if (video) {
      if (Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(hlsUrl);
        hls.attachMedia(video);
        hls.on(Hls.Events.MANIFEST_PARSED, () => {
          video.play();
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = hlsUrl;
        video.addEventListener('loadedmetadata', () => {
          video.play();
        });
      }
    }

    return () => {
      if (video) {
        video.pause();
      }
    };
  }, [hlsUrl]);

  return (
    <div>
      <h1>HLS Video Streaming</h1>
      <video ref={videoRef} controls style={{ width: '100%', maxWidth: '800px' }} />
    </div>
  );
};

