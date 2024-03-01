package token

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xiaowuzai/simplebank/util"
)

func TestJWTMaker(t *testing.T) {
	// 测试JWTMaker的New方法
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.RandomString(6)
	duration := time.Minute
	// 测试JWTMaker的Create方法
	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	// 测试JWTMaker的Verify方法
	payloadV, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payloadV)
	require.NotEmpty(t, payloadV.ID)

	require.Equal(t, payload.Username, payloadV.Username)
	require.WithinDuration(t, payload.ExpiredAt, payloadV.ExpiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, payloadV.IssuedAt, time.Second)
}

func TestVerifyToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), "", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, token)

	payloadV, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), ErrExpiredToken.Error())) // TODO: 如何断言错误？
	require.Nil(t, payloadV)
}
